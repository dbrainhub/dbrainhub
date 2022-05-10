package dbs

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockStatusQuerier struct {
	InitVal int64
	Delta   int64
	Count   int64
}

func (m *MockStatusQuerier) StatementCount(ctx context.Context) (int64, error) {
	m.Count++
	return m.InitVal + m.Count*m.Delta, nil
}
func (m *MockStatusQuerier) TransactionCount(ctx context.Context) (int64, error) {
	m.Count++
	return m.InitVal + m.Count*m.Delta, nil
}

func TestGetQPS(t *testing.T) {
	manager, err := NewDBIndexManager(context.Background(), &MockStatusQuerier{
		InitVal: 100,
		Delta:   5,
		Count:   0,
	})
	assert.Nil(t, err)

	// real_sleep_ms >= 100ms
	time.Sleep(100 * time.Millisecond)

	qps, err := manager.GetQPS(context.Background())
	assert.Nil(t, err)
	assert.Less(t, qps, 100.0)
	assert.Greater(t, qps, 90.0)
}

func TestGetTPS(t *testing.T) {
	manager, err := NewDBIndexManager(context.Background(), &MockStatusQuerier{
		InitVal: 100,
		Delta:   5,
		Count:   0,
	})
	assert.Nil(t, err)

	// real_sleep_ms >= 100ms
	time.Sleep(100 * time.Millisecond)

	qps, err := manager.GetTPS(context.Background())
	assert.Nil(t, err)
	assert.Less(t, qps, 50.0)
	assert.Greater(t, qps, 40.0)
}
