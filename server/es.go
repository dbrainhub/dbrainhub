package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/model/es"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	"github.com/dbrainhub/dbrainhub/utils/rate_limit"
	"github.com/dbrainhub/dbrainhub/utils/search_time"
)

type ESClient interface {
	Send(msg []*es.ESMessage) error
}
type ESClientWithRateLimiter interface {
	ESClient
	rate_limit.RateLimiter
}

func GetDefaultAsyncESClient() ESClient {
	return defaultAsyncEsClient
}

func GetDefaultESClientWithRateLimiter() ESClientWithRateLimiter {
	return defaultEsClientWithRateLimiter
}

func InitDefaultESClient(serverConf *configs.ServerConfig) error {
	defaultEsClientWithRateLimiter = &esClientWithRateLimit{
		ESSender:                 es.NewESSender(es.GetESClient()),
		SlidingWindowRateLimiter: rate_limit.NewSlidingWindowRateLimiter(int64(serverConf.OutputServer.QpsThreshold)),
	}

	defaultAsyncEsClient = &asyncESClient{
		client:    es.NewESSender(es.GetESClient()),
		batchSize: serverConf.OutputServer.ESBatchSize,
		interval:  time.Duration(serverConf.OutputServer.ESIntervalMs) * time.Millisecond,
	}
	go defaultAsyncEsClient.startSendAsync()

	return nil
}

// storage for indices such as cpu/mem
func GetIndicesIndexName() string {
	year, week, _ := search_time.GetYearAndWeek(time.Now().Format("2006-01-02T15:04:05Z"))
	return fmt.Sprintf("instances-%d-%dw", year, week)
}

var defaultAsyncEsClient *asyncESClient
var defaultEsClientWithRateLimiter *esClientWithRateLimit

type esClientWithRateLimit struct {
	es.ESSender
	*rate_limit.SlidingWindowRateLimiter
}

func (e *esClientWithRateLimit) Send(msgs []*es.ESMessage) error {
	return e.ESSender.SendBatch(msgs)
}

type asyncESClient struct {
	client    es.ESSender
	batchSize int
	interval  time.Duration

	sync.Mutex
	buffer []*es.ESMessage
}

func (e *asyncESClient) Send(msg []*es.ESMessage) error {
	e.Lock()
	e.buffer = append(e.buffer, msg...)
	if len(e.buffer) < e.batchSize {
		e.Unlock()
		return nil
	}

	buffer := e.copyBuffer()
	e.Unlock()
	e.sendSync(buffer)
	return nil
}

func (e *asyncESClient) startSendAsync() {
	tick := time.NewTicker(e.interval)

	for {
		select {
		case <-tick.C:
			e.Lock()
			buffer := e.copyBuffer()
			e.Unlock()
			e.sendSync(buffer)
		}
	}
}

// with lock
func (e *asyncESClient) copyBuffer() []*es.ESMessage {
	var res []*es.ESMessage
	res = append(res, e.buffer...)

	e.buffer = e.buffer[:0]
	return res
}

func (e *asyncESClient) sendSync(msgs []*es.ESMessage) {
	if len(msgs) == 0 {
		return
	}
	err := e.client.SendBatch(msgs)
	if err != nil {
		logger.Errorf("es send sync error, err: %v", err)
	}
}
