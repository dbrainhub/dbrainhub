package dbs

import (
	"context"
)

type (
	SlowLogInfoQuerier interface {
		Query(ctx context.Context) (*SlowLogInfo, error)
	}

	SlowLogInfoQuerierFactory interface {
		Create() (SlowLogInfoQuerier, error)
	}

	SlowLogInfo struct {
		IsOpen bool
		Path   string
	}
)
