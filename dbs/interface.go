package dbs

type DBOperationFactory interface {
	CreateVariablesCenter() (VariablesCenter, error)
	CreateVersionQuerier() (DBVersionQuerier, error)
	CreateStatusQuerier() (StatusQuerier, error)
}
