package model

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dbrainhub/dbrainhub/configs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var once = sync.Once{}

func GetDB(ctx context.Context) *gorm.DB {
	once.Do(func() {
		var err error

		c := &configs.GetGlobalConfig().DB
		gormCfg := &gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
		}
		dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
		db, err = gorm.Open(OpenDialectOrPanic(c.Dialect, dsn), gormCfg)
		if err != nil {
			panic(err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.SetMaxIdleConns(64)
		sqlDB.SetMaxOpenConns(128)
		sqlDB.SetConnMaxLifetime(6 * time.Hour)

		// Disalbe timestamp tracking
		db.Callback().Update().Replace("gorm:update_time_stamp", func(db *gorm.DB) {})
	})
	return db
}

func OpenDialectOrPanic(dialect string, dsn string) gorm.Dialector {
	switch dialect {
	case "mysql":
		return mysql.Open(dsn)
	case "postgres":
		return postgres.Open(dsn)
	case "sqlserver":
		return sqlserver.Open(dsn)
	case "sqlite":
		return sqlite.Open(dsn)
	default:
		panic(fmt.Sprintf("invalid dialect: %s", dialect))
	}
}
