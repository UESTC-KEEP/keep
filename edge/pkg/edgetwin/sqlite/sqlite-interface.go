package sqlite

import "time"

type SqliteInterface interface {
	// ConnectToSqlite 使用默认的数据库进行连接
	ConnectToSqlite() error
	// InserBlobIntoMetricsSqlite 向数据库插入数据 主键设置为uuid
	/*
		输入参数：blob接收二进制的数据
	*/
	InserBlobIntoMetricsSqlite(blob []byte) error
	// DeleteTimeOutFromSqlite 从数据库删除已过期的数据
	DeleteTimeOutFromSqlite(ddl time.Time) error
}
