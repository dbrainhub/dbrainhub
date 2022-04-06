package filebeat

import (
	"context"
	"fmt"
	"testing"

	"github.com/dbrainhub/dbrainhub/model"
	"github.com/stretchr/testify/assert"
)

func TestSlowLogPathListener_Listen(t *testing.T) {
	var slowPathArr = []string{"111", "222", "222", "333", "333"}
	listener := NewSlowLogPathListener(&MockSlowLogQuerier{
		slowPathArr: slowPathArr,
	}, 0)

	newSlowPathChan, errChan := listener.Listen(context.Background())
	assert.Equal(t, <-newSlowPathChan, "111")
	assert.Equal(t, <-newSlowPathChan, "222")
	assert.Equal(t, <-newSlowPathChan, "333")
	assert.Equal(t, <-newSlowPathChan, "111")
	assert.Equal(t, len(errChan), 0)
}

func TestSlowLogPathListener_ListenError(t *testing.T) {
	var slowPathArr = []string{"111", "222", "222", "333", "333"}
	listener := NewSlowLogPathListener(&MockSlowLogQuerier{
		slowPathArr: slowPathArr,
		isReturnErr: true,
	}, 0)

	newSlowPathChan, errChan := listener.Listen(context.Background())
	assert.Equal(t, (<-errChan).Error(), "111")
	assert.Equal(t, (<-errChan).Error(), "222")
	assert.Equal(t, (<-errChan).Error(), "222")
	assert.Equal(t, (<-errChan).Error(), "333")
	assert.Equal(t, (<-errChan).Error(), "333")
	assert.Equal(t, len(newSlowPathChan), 0)
}

type MockSlowLogQuerier struct {
	slowPathArr []string
	i           int
	isReturnErr bool
}

func (m *MockSlowLogQuerier) Query(ctx context.Context) (*model.SlowLogInfo, error) {
	var res = &model.SlowLogInfo{
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
