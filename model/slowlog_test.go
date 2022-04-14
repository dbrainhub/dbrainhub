package model

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

const (
	user          = "root"
	passwd        = "123"
	addr          = "127.0.0.1:3306"
	isOpenSlowLog = false
	slowLogPath   = "/usr/local/var/mysql/lypdeMacBook-Pro-slow.log"
	testSwitch    = false
)

// 只用于测试使用，不同机器参数和结果不同，暂时置为不执行。
func TestQuerySlowLog(t *testing.T) {
	if !testSwitch {
		return
	}

	db, _ := newDB(t)
	defer db.Close()

	slowLogInfo, err := NewSlowLogInfoQuerier(db).Query(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, slowLogInfo.IsOpen, isOpenSlowLog)
	assert.Equal(t, slowLogInfo.Path, slowLogPath)
}

func newDB(t *testing.T) (*sql.DB, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/", user, passwd, addr))
	assert.Nil(t, err)
	return db, err
}
