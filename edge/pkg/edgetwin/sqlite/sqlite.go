package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgetwin/config"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/loggerv1.0.1"
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

func ReceiveFromBeehiveAndInsert() {
	cli := NewSqliteCli()
	go func(*Sqlite) {
		for {
			select {
			case <-beehiveContext.Done():
				return
			default:
			}
			ReceiveEdgeTwinMsg(cli)
		}
	}(cli)

}

func ReceiveEdgeTwinMsg(cli *Sqlite) {
	msg, err := beehiveContext.Receive(modules.EdgeTwinModule)
	if err != nil {
		logger.Error(err)
		time.Sleep(1 * time.Second)
	} else {
		logger.Trace("接收消息 msg: ", msg)

		// 提高速度
		go func() {
			err := cli.InserBlobIntoMetricsSqlite(msg.Content.([]byte))
			if err != nil {
				logger.Error(err)
				return
			}
			resp := msg.NewRespByMessage(&msg, " message received ")
			beehiveContext.SendResp(*resp)
		}()
	}
}
