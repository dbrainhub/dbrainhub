package mysql

import (
	"context"
	"testing"

	"github.com/dbrainhub/dbrainhub/dbs"
	"github.com/stretchr/testify/assert"
)

func newDBInfo() *dbs.DBInfo {
	return &dbs.DBInfo{
		IP:     "127.0.0.1",
		Port:   3306,
		User:   "root",
		Passwd: "123",
	}
}

func TestGetTransactions(t *testing.T) {
	factory := NewMysqlOperationFactory(newDBInfo())
	querier, err := factory.CreateStatusCenter()
	assert.Nil(t, err)
	cnt, err := querier.TransactionCount(context.Background())
	assert.Nil(t, err)
	assert.NotEqual(t, cnt, 0)

	cnt, err = querier.StatementCount(context.Background())
	assert.Nil(t, err)
	assert.NotEqual(t, cnt, 0)
}
