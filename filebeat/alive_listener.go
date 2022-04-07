package filebeat

import (
	"context"
	"fmt"
	"time"

	"github.com/dbrainhub/dbrainhub/utils"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

type (
	AliveListener interface {
		Listen(ctx context.Context)
	}

	AliveErrorCallback func(ctx context.Context, err error)
)

func NewAliveListener(host string, client utils.HttpClient, interval, delay time.Duration, callback AliveErrorCallback) AliveListener {
	return &filebeatAliveListener{
		host:          host,
		client:        client,
		interval:      interval,
		delay:         delay,
		errorCallback: callback,
	}
}

// alive check details: https://www.elastic.co/guide/en/beats/filebeat/current/http-endpoint.html
// `http.enabled == true` is guaranteed in conf validater.
type filebeatAliveListener struct {
	host          string
	client        utils.HttpClient
	interval      time.Duration
	delay         time.Duration // some time for filebeat startup
	errorCallback AliveErrorCallback
}

func (f *filebeatAliveListener) Listen(ctx context.Context) {
	url := fmt.Sprintf("http://%s/?pretty", f.host)

	go func() {
		time.Sleep(f.delay)

		for {
			resp, err := f.client.Send(ctx, url, "GET", "")
			logger.Infof("alive check resp: %v, err:%v", resp, err)
			if err != nil {
				f.errorCallback(ctx, err)
			}

			time.Sleep(f.interval)
		}
	}()
}
