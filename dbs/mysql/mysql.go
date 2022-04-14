package mysql

import "github.com/dbrainhub/dbrainhub/dbs"

const (
	MysqlType = "mysql"
)

func NewMysqlOperationFactory(dbInfo *dbs.DBInfo) dbs.DBOperationFactory {
	return &mysqlOperationFactory{
		dbInfo: dbInfo,
	}
}

type mysqlOperationFactory struct {
	dbInfo *dbs.DBInfo
}

func (m *mysqlOperationFactory) CreateSlowlogQuerier() (dbs.SlowLogInfoQuerier, error) {
	db, err := m.dbInfo.GetSQLDB(MysqlType)
	if err != nil {
		return nil, err
	}
	return &mysqlSlowLogInfoQuerier{
		db: db,
	}, nil
}

func (m *mysqlOperationFactory) CreateVersionQuerier() (dbs.DBVersionQuerier, error) {
	db, err := m.dbInfo.GetSQLDB(MysqlType)
	if err != nil {
		return nil, err
	}
	return &mysqlDBVersion{
		db: db,
	}, nil
}
