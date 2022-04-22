package dbs

import (
	"context"
	"fmt"
	"time"
)

const (
	InvalidQPS = -1
	InvalidTPS = -1
)

type (
	DBIndexManager interface {
		GetQPS(ctx context.Context) (float64, error)
		GetTPS(ctx context.Context) (float64, error)
	}

	StatusQuerier interface {
		QueryStatementCount(ctx context.Context) (int64, error)
		QueryTransactionCount(ctx context.Context) (int64, error)
	}
)

func NewDBIndexManager(ctx context.Context, statusQuerier StatusQuerier) (DBIndexManager, error) {
	res := &defaultIndexManager{
		statusQuerier: statusQuerier,
	}
	_, _ = res.GetQPS(ctx)
	_, _ = res.GetTPS(ctx) // need one more calls
	return res, nil
}

type defaultIndexManager struct {
	statusQuerier StatusQuerier

	lastStatementCount int64
	lastStatementNs    int64

	lastTransactionCount int64
	lastTransactionNs    int64
}

// Note: The interval that call `NewDBIndexManager` and `GetQPS` is too short may cause return error
func (d *defaultIndexManager) GetQPS(ctx context.Context) (float64, error) {
	count, err := d.statusQuerier.QueryStatementCount(ctx)
	if err != nil {
		return InvalidQPS, err
	}

	currentNs := time.Now().UnixNano()
	timeDelta := currentNs - d.lastStatementNs
	if timeDelta == 0 {
		return InvalidQPS, fmt.Errorf("GetQPS time_delta is 0")
	}

	res := float64((count - d.lastStatementCount)) / (float64(timeDelta) / 1e9)
	d.lastStatementCount = count
	d.lastStatementNs = currentNs

	if res < 0 {
		return InvalidQPS, fmt.Errorf("GetQPS conut_delta is negative")
	}
	return res, nil
}

// Note: The interval that call `NewDBIndexManager` and `GetTPS` is too short may cause return error
func (d *defaultIndexManager) GetTPS(ctx context.Context) (float64, error) {
	count, err := d.statusQuerier.QueryTransactionCount(ctx)
	if err != nil {
		return InvalidTPS, err
	}

	currentNs := time.Now().UnixNano()
	timeDelta := currentNs - d.lastTransactionNs
	if timeDelta == 0 {
		return InvalidQPS, fmt.Errorf("GetTPS time_delta is 0")
	}

	res := float64((count - d.lastTransactionCount)) / float64(timeDelta) * 1e9
	d.lastTransactionCount = count
	d.lastTransactionNs = currentNs

	if res < 0 {
		return InvalidTPS, fmt.Errorf("GetTPS conut_delta is negative")
	}
	return res, nil
}
