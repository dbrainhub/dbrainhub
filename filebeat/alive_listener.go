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

	AliveSuccCallback  func(ctx context.Context, resp string)
	AliveErrorCallback func(ctx context.Context, err error)

	AliveListenerCallback struct {
		ErrorCallback AliveErrorCallback
		SuccCallback  AliveSuccCallback
	}
)

func NewAliveListener(host string, client utils.HttpClient, interval time.Duration, callbacks *AliveListenerCallback) AliveListener {
	return &filebeatAliveListener{
		host:      host,
		client:    client,
		interval:  interval,
		callbacks: callbacks,
	}
}

// alive check details: https://www.elastic.co/guide/en/beats/filebeat/current/http-endpoint.html
// `http.enabled == true` is guaranteed in conf validater.
type filebeatAliveListener struct {
	host      string
	client    utils.HttpClient
	interval  time.Duration
	callbacks *AliveListenerCallback
}

func (f *filebeatAliveListener) Listen(ctx context.Context) {
	url := fmt.Sprintf("http://%s/?pretty", f.host)

	go func() {
		for {
			resp, err := f.client.Send(ctx, url, "GET", "")
			logger.Infof("alive check resp: %s, err:%v", string(resp), err)
			if err != nil {
				f.callbacks.ErrorCallback(ctx, err)
			} else {
				f.callbacks.SuccCallback(ctx, string(resp))
			}

			time.Sleep(f.interval)
		}
	}()
}
