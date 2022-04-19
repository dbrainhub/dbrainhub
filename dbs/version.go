package dbs

import (
	"context"
)

type (
	DBVersionQuerier interface {
		Query(ctx context.Context) (*DBVersion, error)
	}
	DBVersion struct {
		Version string
	}
)
