package server

import (
	"sync"
	"time"

	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

type ESClientAsync interface {
	// send async and will logger error when send failed
	Send(msg *model.ESMessage)
}

var defaultEsClientAsync ESClientAsync

func GetDefaultEsClientAsync() ESClientAsync {
	return defaultEsClientAsync
}

func InitDefaultEsClientAsync(serverConf *configs.ServerConfig) error {
	client, err := model.NewEsClient(serverConf.OutputServer.EsAddresses)
	if err != nil {
		return err
	}

	esClientAsync := &esClientAsync{
		client:    client,
		batchSize: serverConf.OutputServer.ESBatchSize,
		interval:  time.Duration(serverConf.OutputServer.ESIntervalMs) * time.Millisecond,
	}
	defaultEsClientAsync = esClientAsync

	go esClientAsync.sendAsync()
	return nil
}

type esClientAsync struct {
	client    model.EsSender
	batchSize int
	interval  time.Duration

	sync.Mutex
	buffer []*model.ESMessage
}

func (e *esClientAsync) Send(msg *model.ESMessage) {
	e.Lock()
	e.buffer = append(e.buffer, msg)
	if len(e.buffer) < e.batchSize {
		e.Unlock()
		return
	}

	buffer := e.copyBuffer()
	e.Unlock()
	e.sendSync(buffer)
}

func (e *esClientAsync) sendAsync() {
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
func (e *esClientAsync) copyBuffer() []*model.ESMessage {
	var res []*model.ESMessage
	res = append(res, e.buffer...)

	e.buffer = e.buffer[:0]
	return res
}

func (e *esClientAsync) sendSync(msgs []*model.ESMessage) {
	if len(msgs) == 0 {
		return
	}
	err := e.client.SendBatch(msgs)
	if err != nil {
		logger.Errorf("es send sync error, err: %v", err)
	}
}
