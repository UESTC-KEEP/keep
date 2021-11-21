package healthzhub

type HealthzHub interface {
	// InsertIntoSqlite 向sqlite插入数据
	InsertIntoSqlite(blod []byte) error
}
