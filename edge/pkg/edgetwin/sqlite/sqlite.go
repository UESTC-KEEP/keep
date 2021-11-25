package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"github.com/wonderivan/logger"
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgetwin/config"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/kplogger"
	"time"
)

type Sqlite struct {
}

var conn *sql.DB

func (sq *Sqlite) ConnectToSqlite() error {
	db, err := sql.Open("sqlite3", config.Config.SqliteFilePath)
	conn = db
	if err != nil {
		logger.Error("Failed to open sqlite", err)
	}
	return err
}

func (sq *Sqlite) InserBlobIntoMetricsSqlite(blob []byte) error {
	//插入数据
	if conn == nil {
		new(Sqlite).ConnectToSqlite()
	}
	stmt, err := conn.Prepare("INSERT INTO metrics(uuid,data) values(?,?)")
	if err != nil {
		logger.Error(err)
	}

	res, err := stmt.Exec((uuid.NewV4()).String(), blob)
	if err != nil {
		logger.Error(err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		logger.Error(err)
	}
	return err
}

func NewSqliteCli() *Sqlite {
	return new(Sqlite)
}

var ListenBeehiveChannel = make(chan bool)

func ReceiveFromBeehiveAndInsert() {
	cli := NewSqliteCli()
	go func() {
		for {
			//logger.Error("---------------")
			select {
			case <-ListenBeehiveChannel:
				// 接到退出信号 即停止收消息
				logger.Warn("收到退出信号....")
				return
			//设置超时时间 超时就不管了
			case <-time.After(config.Config.BeehiveTimeout):
				msg, err := beehiveContext.Receive(modules.EdgeTwinModule)
				if err != nil {
					kplogger.Error(err)
					time.Sleep(5 * time.Second)
				} else {
					logger.Trace("接收消息 msg: ", msg)
					if msg.Content == nil {
						logger.Warn("消息为空....")
					} else {
						resp := msg.NewRespByMessage(&msg, " message received ")
						beehiveContext.SendResp(*resp)
						err := cli.InserBlobIntoMetricsSqlite(msg.Content.([]byte))
						if err != nil {
							logger.Error(err)
						}
					}
				}
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

}
