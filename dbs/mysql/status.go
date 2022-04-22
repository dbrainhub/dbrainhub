package mysql

import (
	"context"
	"database/sql"
	"fmt"
)

type mysqlStatusCenter struct {
	db *sql.DB
}

func (m *mysqlStatusCenter) StatementCount(ctx context.Context) (int64, error) {
	return m.queryQuertions(ctx)
}

func (m *mysqlStatusCenter) TransactionCount(ctx context.Context) (int64, error) {
	var res int64
	selectCount, err := m.queryComSelect(ctx)
	if err != nil {
		return 0, err
	}
	res += selectCount

	insertCount, err := m.queryComInsert(ctx)
	if err != nil {
		return 0, err
	}
	res += insertCount

	updateCount, err := m.queryComUpdate(ctx)
	if err != nil {
		return 0, err
	}
	res += updateCount

	deleteCount, err := m.queryComDelete(ctx)
	if err != nil {
		return 0, err
	}
	res += deleteCount

	commitCount, err := m.queryComCommit(ctx)
	if err != nil {
		return 0, err
	}
	res += commitCount

	rollbackCount, err := m.queryComRollback(ctx)
	if err != nil {
		return 0, err
	}
	res += rollbackCount

	updateMultiCount, err := m.queryComUpdateMulti(ctx)
	if err != nil {
		return 0, err
	}
	res += updateMultiCount

	insertSelectCount, err := m.queryComInsertSelect(ctx)
	if err != nil {
		return 0, err
	}
	res += insertSelectCount

	deleteMultiCount, err := m.queryComDeleteMulti(ctx)
	if err != nil {
		return 0, err
	}
	res += deleteMultiCount

	return res, nil
}

// The difference between `Questions` and `Queries` can  refer to https://dev.mysql.com/doc/refman/8.0/en/server-status-variables.html
func (m *mysqlStatusCenter) queryQuertions(ctx context.Context) (int64, error) {
	const variableName = "Questions"
	var res int64
	if err := m.query(ctx, variableName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mysqlStatusCenter) queryComSelect(ctx context.Context) (int64, error) {
	const variableName = "Com_select"
	var res int64
	if err := m.query(ctx, variableName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mysqlStatusCenter) queryComInsert(ctx context.Context) (int64, error) {
	const variableName = "Com_insert"
	var res int64
	if err := m.query(ctx, variableName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mysqlStatusCenter) queryComUpdate(ctx context.Context) (int64, error) {
	const variableName = "Com_update"
	var res int64
	if err := m.query(ctx, variableName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mysqlStatusCenter) queryComDelete(ctx context.Context) (int64, error) {
	const variableName = "Com_delete"
	var res int64
	if err := m.query(ctx, variableName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mysqlStatusCenter) queryComInsertSelect(ctx context.Context) (int64, error) {
	const variableName = "Com_insert_select"
	var res int64
	if err := m.query(ctx, variableName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

// update multi tables
func (m *mysqlStatusCenter) queryComUpdateMulti(ctx context.Context) (int64, error) {
	const variableName = "Com_update_multi"
	var res int64
	if err := m.query(ctx, variableName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mysqlStatusCenter) queryComDeleteMulti(ctx context.Context) (int64, error) {
	const variableName = "Com_delete_multi"
	var res int64
	if err := m.query(ctx, variableName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

// note: there is a update statement in a transaction: Com_commit++, Com_update not changed
func (m *mysqlStatusCenter) queryComCommit(ctx context.Context) (int64, error) {
	const variableName = "Com_commit"
	var res int64
	if err := m.query(ctx, variableName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mysqlStatusCenter) queryComRollback(ctx context.Context) (int64, error) {
	const variableName = "Com_rollback"
	var res int64
	if err := m.query(ctx, variableName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mysqlStatusCenter) query(ctx context.Context, name string, value interface{}) error {
	rows, err := m.db.Query(fmt.Sprintf(`show global status like '%s'`, name))
	if err != nil {
		return err
	}
	defer rows.Close()

	var variableName string
	for rows.Next() {
		if err = rows.Scan(&variableName, value); err != nil {
			return err
		}
		break
	}
	return nil
}
