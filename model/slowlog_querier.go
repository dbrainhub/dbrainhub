package model

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type (
	SlowLogInfoQuerier interface {
		Query(ctx context.Context) (*SlowLogInfo, error)
	}

	SlowLogInfo struct {
		IsOpen bool
		Path   string
	}
)

func NewSlowLogInfoQuerier(db *sql.DB) SlowLogInfoQuerier {
	return &mysqlSlowLogInfoQuerier{
		db: db,
	}
}

type mysqlSlowLogInfoQuerier struct {
	db *sql.DB
}

func (m *mysqlSlowLogInfoQuerier) Query(ctx context.Context) (*SlowLogInfo, error) {
	var res SlowLogInfo
	var err error
	res.IsOpen, err = m.querySlowLogIsOpen(ctx)
	if err != nil {
		return nil, err
	}
	res.Path, err = m.querySlowLogPath(ctx)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mysqlSlowLogInfoQuerier) querySlowLogIsOpen(ctx context.Context) (bool, error) {
	rows, err := m.db.Query(`show global variables like 'slow_query_log'`)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var variableName, value string
	for rows.Next() {
		if err = rows.Scan(&variableName, &value); err != nil {
			return false, err
		}
		break
	}
	return value == "ON", nil
}

func (m *mysqlSlowLogInfoQuerier) querySlowLogPath(ctx context.Context) (string, error) {
	rows, err := m.db.Query(`show global variables like 'slow_query_log_file'`)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var variableName, value string
	for rows.Next() {
		if err = rows.Scan(&variableName, &value); err != nil {
			return "", err
		}
		break
	}
	return value, nil
}
