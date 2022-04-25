package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dbrainhub/dbrainhub/dbs"
)

type mysqlVariablesQuerier struct {
	db *sql.DB
}

func (m *mysqlVariablesQuerier) SlowlogInfo(ctx context.Context) (*dbs.SlowLogInfo, error) {
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

func (m *mysqlVariablesQuerier) DataDir(ctx context.Context) (string, error) {
	const DataDirVar = "datadir"

	var res string
	if err := m.queryGlobal(ctx, DataDirVar, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (m *mysqlVariablesQuerier) querySlowLogIsOpen(ctx context.Context) (bool, error) {
	const SlowQueryLogVar = "slow_query_log"

	var isOpen string
	if err := m.queryGlobal(ctx, SlowQueryLogVar, &isOpen); err != nil {
		return false, err
	}

	return isOpen == "ON", nil
}

func (m *mysqlVariablesQuerier) querySlowLogPath(ctx context.Context) (string, error) {
	const SlowlogPathVar = "slow_query_log_file"

	var res string
	if err := m.queryGlobal(ctx, SlowlogPathVar, &res); err != nil {
		return "", err
	}

	return res, nil
}

func (m *mysqlVariablesQuerier) queryGlobal(ctx context.Context, varNameLike string, val interface{}) error {
	rows, err := m.db.Query(fmt.Sprintf("show global variables like '%s'", varNameLike))
	if err != nil {
		return err
	}
	defer rows.Close()

	var variableName string
	for rows.Next() {
		if err = rows.Scan(&variableName, val); err != nil {
			return err
		}
		break
	}
	return nil
}
