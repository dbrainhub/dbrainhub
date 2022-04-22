package dbs

import (
	"context"
)

type (
	VariablesCenter interface {
		SlowlogInfo(ctx context.Context) (*SlowLogInfo, error)
		DataDir(ctx context.Context) (string, error)
	}

	SlowLogInfo struct {
		IsOpen bool
		Path   string
	}
)
