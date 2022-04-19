package filebeat

import (
	"context"
	"fmt"
	"testing"

	"github.com/dbrainhub/dbrainhub/dbs"
	"github.com/stretchr/testify/assert"
)

func TestSlowLogPathListener_Listen(t *testing.T) {
	var slowPathArr = []string{"111", "222", "222", "333", "333"}

	// use chan for testing logic in another gorounte
	msgChan := make(chan string)
	callbacks := &SlowLogPathCallback{
		ChangedCallback: func(ctx context.Context, oldPath, newPath string) {
			msgChan <- newPath
		},
		ErrorCallback: nil,
	}

	listener := NewSlowLogPathListener(&MockSlowLogQuerier{
		slowPathArr: slowPathArr,
	}, 0, callbacks)

	listener.Listen(context.Background())
	assert.Equal(t, <-msgChan, "111")
	assert.Equal(t, <-msgChan, "222")
	assert.Equal(t, <-msgChan, "333")
	assert.Equal(t, <-msgChan, "111")
}

func TestSlowLogPathListener_ListenError(t *testing.T) {
	var slowPathArr = []string{"111", "222", "222", "333", "333"}

	msgChan := make(chan string)
	callbacks := &SlowLogPathCallback{
		ChangedCallback: nil,
		ErrorCallback: func(ctx context.Context, err error) {
			msgChan <- err.Error()
		},
	}

	listener := NewSlowLogPathListener(&MockSlowLogQuerier{
		slowPathArr: slowPathArr,
		isReturnErr: true,
	}, 0, callbacks)

	listener.Listen(context.Background())
	assert.Equal(t, <-msgChan, "111")
	assert.Equal(t, <-msgChan, "222")
	assert.Equal(t, <-msgChan, "222")
	assert.Equal(t, <-msgChan, "333")
	assert.Equal(t, <-msgChan, "333")

}

type MockSlowLogQuerier struct {
	slowPathArr []string
	i           int
	isReturnErr bool
}

func (m *MockSlowLogQuerier) Query(ctx context.Context) (*dbs.SlowLogInfo, error) {
	var res = &dbs.SlowLogInfo{
		Path:   m.slowPathArr[m.i%len(m.slowPathArr)],
		IsOpen: true,
	}
	err := fmt.Errorf(m.slowPathArr[m.i%len(m.slowPathArr)])
	m.i++

	if m.isReturnErr {
		return nil, err
	}
	return res, nil
}
