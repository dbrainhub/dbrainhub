package mysql

import (
	"context"
	"database/sql"

	"github.com/dbrainhub/dbrainhub/dbs"
)

type mysqlSlowLogInfoQuerier struct {
	db *sql.DB
}

func (m *mysqlSlowLogInfoQuerier) Query(ctx context.Context) (*dbs.SlowLogInfo, error) {
	var res dbs.SlowLogInfo
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
