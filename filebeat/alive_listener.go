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
		Listen(ctx context.Context) chan error
	}
)

func NewAliveListener(host string, client utils.HttpClient, duration, delay time.Duration) AliveListener {
	return &filebeatAliveListener{
		host:     host,
		client:   client,
		duration: duration,
		delay:    delay,
	}
}

type filebeatAliveListener struct {
	host     string
	client   utils.HttpClient
	duration time.Duration
	delay    time.Duration
}

func (f *filebeatAliveListener) Listen(ctx context.Context) chan error {
	var errCh = make(chan error, 2)
	url := fmt.Sprintf("http://%s/?pretty", f.host)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("filebeat alive listener panic, err:%v", err) // TODO: PANIC?
			}
		}()
		time.Sleep(f.delay)
		for {
			resp, err := f.client.Send(ctx, url, "GET", "")
			logger.Infof("alive check resp: %v, err:%v", resp, err)
			if err != nil {
				errCh <- err
			}

			time.Sleep(f.duration)
		}
	}()
	return errCh
}
