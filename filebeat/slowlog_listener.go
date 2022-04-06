package filebeat

import (
	"context"
	"time"

	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

type (
	SlowLogPathListener interface {
		Listen(ctx context.Context) (chan string, chan error)
	}
)

func NewSlowLogPathListener(slowLogQuerier model.SlowLogInfoQuerier, duration time.Duration) SlowLogPathListener {
	return &slowLogPathListener{
		slowLogQuerier: slowLogQuerier,
		duration:       duration,
	}
}

type slowLogPathListener struct {
	slowLogQuerier model.SlowLogInfoQuerier
	duration       time.Duration
}

func (s *slowLogPathListener) Listen(ctx context.Context) (chan string, chan error) {
	var newPathChan = make(chan string, 2)
	var errChan = make(chan error, 2)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("slowlog path listener panic, err:%v", err) // TODO: PANIC?
			}
		}()

		var slowLogPath string
		for {
			slowLogInfo, err := s.slowLogQuerier.Query(ctx)
			if err != nil {
				logger.Warnf("listen slowlog path error, err: %v", err)
				errChan <- err
			} else if slowLogInfo.Path != slowLogPath {
				logger.Infof("slowlog path update, old: %s, new: %s", slowLogPath, slowLogInfo.Path)
				slowLogPath = slowLogInfo.Path
				newPathChan <- slowLogPath
			}

			time.Sleep(s.duration)
		}
	}()
	return newPathChan, errChan
}
