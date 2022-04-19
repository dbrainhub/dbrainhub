package dbs

type DBOperationFactory interface {
	CreateSlowlogQuerier() (SlowLogInfoQuerier, error)
	CreateVersionQuerier() (DBVersionQuerier, error)
}
