package mysql

import (
	"context"
	"database/sql"

	"github.com/dbrainhub/dbrainhub/dbs"
)

type mysqlDBVersion struct {
	db *sql.DB
}

func (m *mysqlDBVersion) Query(ctx context.Context) (*dbs.DBVersion, error) {
	var res dbs.DBVersion

	rows, err := m.db.Query(`select version() as version`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&res.Version); err != nil {
			return nil, err
		}
		break
	}
	return &res, nil
}
