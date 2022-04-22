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
	querier, err := factory.CreateStatusQuerier()
	assert.Nil(t, err)
	cnt, err := querier.QueryTransactionCount(context.Background())
	assert.Nil(t, err)
	assert.NotEqual(t, cnt, 0)

	cnt, err = querier.QueryStatementCount(context.Background())
	assert.Nil(t, err)
	assert.NotEqual(t, cnt, 0)
}
