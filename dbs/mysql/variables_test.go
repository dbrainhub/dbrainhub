package mysql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDataDir(t *testing.T) {
	factory := NewMysqlOperationFactory(newDBInfo())
	vs, err := factory.CreateVariablesCenter()
	assert.Nil(t, err)

	dataDir, err := vs.DataDir(context.Background())
	assert.Nil(t, err)
	assert.NotEmpty(t, dataDir, "")
}

func TestGetSlowlogInfo(t *testing.T) {
	factory := NewMysqlOperationFactory(newDBInfo())
	vs, err := factory.CreateVariablesCenter()
	assert.Nil(t, err)

	slowlogInfo, err := vs.SlowlogInfo(context.Background())
	assert.Nil(t, err)
	assert.NotEqual(t, slowlogInfo.Path, "")
}
