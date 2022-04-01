package dao

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

const (
	User          = "root"
	Passwd        = "123"
	Addr          = "127.0.0.1:3306"
	IsOpenSlowLog = false
	SlowLogPath   = "/usr/local/var/mysql/lypdeMacBook-Pro-slow.log"
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
	assert.Equal(t, slowLogInfo.IsOpen, IsOpenSlowLog)
	assert.Equal(t, slowLogInfo.Path, SlowLogPath)
}

func newDB(t *testing.T) (*sql.DB, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/", User, Passwd, Addr))
	assert.Nil(t, err)
	return db, err
}
