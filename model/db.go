package model

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dbrainhub/dbrainhub/configs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var once = sync.Once{}

func GetDB(ctx context.Context) *gorm.DB {
	once.Do(func() {
		var err error
		c := &configs.GetGlobalServerConfig().DB
		url := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
		db, err = gorm.Open(c.Dialect, url)
		if err != nil {
			panic(err)
		}

		db.DB().SetMaxIdleConns(64)
		db.DB().SetMaxOpenConns(128)
		db.DB().SetConnMaxLifetime(6 * time.Hour)

		// Disable table name's pluralization for GORM
		db.SingularTable(true)
		// Disalbe timestamp tracking
		db.Callback().Update().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {})
	})
	return db
}
