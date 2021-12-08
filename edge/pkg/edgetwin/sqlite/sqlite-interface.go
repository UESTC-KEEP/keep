package sqlite

type SqliteInterface interface {
	// ConnectToSqlite 使用默认的数据库进行连接
	ConnectToSqlite() error
	// InserBlobIntoMetricsSqlite 向数据库插入数据 主键设置为uuid
	/*
		输入参数：blob接收二进制的数据
	*/
	InserBlobIntoMetricsSqlite(content []byte , time string) error
	// DeleteTimeOutFromSqlite 从数据库删除已过期的数据
	DeleteTimeOutFromSqlite(ddl int64) error

	// SelectTimeFromSqliteToCloud 检索断连后的数据并同步给云端
	SelectTimeFromSqliteToCloud( begintime int64) error
}
