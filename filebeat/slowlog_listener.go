package filebeat

import (
	"context"
	"time"

	"github.com/dbrainhub/dbrainhub/dbs"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

type (
	SlowLogPathListener interface {
		Listen(ctx context.Context)
	}

	SlowLogPathChangedCallback func(ctx context.Context, oldPath, newPath string)
	SlowlogErrorCallback       func(ctx context.Context, err error)
	SlowLogPathCallback        struct {
		ChangedCallback SlowLogPathChangedCallback
		ErrorCallback   SlowlogErrorCallback
	}
)

func NewSlowLogPathListener(slowLogQuerier slowlogQuerier, interval time.Duration, callbacks *SlowLogPathCallback) SlowLogPathListener {
	return &slowLogPathListener{
		slowLogQuerier: slowLogQuerier,
		interval:       interval,
		callbacks:      callbacks,
	}
}

type slowlogQuerier interface {
	SlowlogInfo(ctx context.Context) (*dbs.SlowLogInfo, error)
}

type slowLogPathListener struct {
	slowLogQuerier slowlogQuerier
	interval       time.Duration
	callbacks      *SlowLogPathCallback
}

func (s *slowLogPathListener) Listen(ctx context.Context) {
	go func() {
		var slowLogPath string
		for {
			slowLogInfo, err := s.slowLogQuerier.SlowlogInfo(ctx)
			if err != nil {
				logger.Warnf("listen slowlog path error, err: %v", err)

				s.callbacks.ErrorCallback(ctx, err)
			} else if slowLogInfo.Path != slowLogPath {
				logger.Infof("slowlog path update, old: %s, new: %s", slowLogPath, slowLogInfo.Path)

				s.callbacks.ChangedCallback(ctx, slowLogPath, slowLogInfo.Path)
				slowLogPath = slowLogInfo.Path
			} else {
				logger.Infof("slowlog path unchange.")
			}

			time.Sleep(s.interval)
		}
	}()
}
