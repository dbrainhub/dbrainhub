package dbs

import (
	"database/sql"
	"fmt"
)

type DBInfo struct {
	IP     string
	Port   int
	User   string
	Passwd string

	sharedDB *sql.DB
}

func (d *DBInfo) GetSQLDB(dbtype string) (*sql.DB, error) {
	if d.sharedDB != nil {
		return d.sharedDB, nil
	}

	var err error
	d.sharedDB, err = d.NewSQLDB(dbtype)
	return d.sharedDB, err
}

func (d *DBInfo) NewSQLDB(dbtype string) (*sql.DB, error) {
	return sql.Open(dbtype,
		fmt.Sprintf("%s:%s@tcp(%s:%d)/", d.User, d.Passwd, d.IP, d.Port))
}
