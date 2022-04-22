package mysql

import (
	"github.com/dbrainhub/dbrainhub/dbs"
	_ "github.com/go-sql-driver/mysql"
)

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

func (m *mysqlOperationFactory) CreateVariablesCenter() (dbs.VariablesCenter, error) {
	db, err := m.dbInfo.GetSQLDB(MysqlType)
	if err != nil {
		return nil, err
	}
	return &mysqlVariablesQuerier{
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

func (m *mysqlOperationFactory) CreateStatusQuerier() (dbs.StatusQuerier, error) {
	db, err := m.dbInfo.GetSQLDB(MysqlType)
	if err != nil {
		return nil, err
	}
	return &mysqlStatusQuerier{
		db: db,
	}, nil
}
