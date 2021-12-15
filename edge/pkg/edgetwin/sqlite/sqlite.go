package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgetwin/config"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
	logger "keep/pkg/util/loggerv1.0.1"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	//"github.com/robfig/cron/v3"
)

type Sqlite struct {
	conn *sql.DB
}

//var conn *sql.DB

func (sq *Sqlite) ConnectToSqlite() error {
	var err error
	sq.conn, err = sql.Open("sqlite3", config.Config.SqliteFilePath)
	if err != nil {
		logger.Error("Failed to open sqlite", err)
	}
	return err
}

func (sq *Sqlite) InserBlobIntoMetricsSqlite(content []byte, time string) error {
	//插入数据
	if sq.conn == nil {
		err := sq.ConnectToSqlite()
		if err != nil {
			return err
		}
	}
	stmt, err := sq.conn.Prepare("INSERT INTO logedgeagent(time , content) values(?,?)")
	if err != nil {
		logger.Error(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(time, string(content))
	if err != nil {
		logger.Error(err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (sq *Sqlite) DeleteTimeOutFromSqlite(ddl int64) error {
	stmt, err := sq.conn.Prepare("DELETE FROM logedgeagent WHERE time <= ?")
	if err != nil {
		logger.Error(err)
		return err
	}
	_, err = stmt.Exec(ddl)
	if err != nil {
		logger.Error(err)
		return err
	}
	// affect, err = res.RowsAffected()
	return err
}

func (sq *Sqlite) SelectTimeFromSqliteToCloud(begintime int64) error {

	stmt, err := sq.conn.Prepare("SELECT content FROM logedgeagent where time > ?")
	if err != nil {
		logger.Error(err)
		return err
	}

	rows, err := stmt.Query(begintime)
	if err != nil {
		logger.Error(err)
		return err
	}
	for rows.Next() {
		var content []byte
		err := rows.Scan(&content)
		if err != nil {
			logger.Error(err)
			return err
		}
		var msg model.Message
		err = json.Unmarshal(content, &msg)
		if err != nil {
			logger.Error(err)
			return err
		}
		msg.Router.Source = modules.EdgeTwinModule
		beehiveContext.Send(modules.EdgePublisherModule, msg)
	}

	return err
}

func NewSqliteCli() *Sqlite {
	conn, err := sql.Open("sqlite3", config.Config.SqliteFilePath)
	if err != nil {
		logger.Error("Failed to open sqlite", err)
	}
	conn.Exec("CREATE TABLE IF NOT EXISTS logedgeagent(time TEXT ,content TEXT)")
	sq := &Sqlite{
		conn: conn,
	}
	//cron := DeletePeriod(sq)
	//defer cron.Stop()
	return sq
}

func ReceiveFromBeehiveAndInsert() {
	cli := NewSqliteCli()
	go func(*Sqlite) {
		c := DeletePeriod(cli)
		defer c.Stop()
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
	msg, err := beehiveContext.Receive(modules.EdgeTwinGroup)
	if err != nil {
		logger.Error(err)
		time.Sleep(1 * time.Second)
	} else {
		//logger.Trace("接收消息 msg: ", msg)
		fmt.Printf("edtwin 接收消息 msg:%#v ", msg.GetID())
		// 提高速度
		go func() {
			content, err := json.Marshal(msg)
			if err != nil {
				logger.Error(err)
				return
			}

			err1 := cli.InserBlobIntoMetricsSqlite(content, strconv.FormatInt(msg.Header.Timestamp, 10))
			if err1 != nil {
				logger.Error(err1)
				return
			}
		}()
	}
}

func DeletePeriod(sq *Sqlite) *cron.Cron {
	c := cron.New()
	c.AddFunc("@every 48h", func() {
		err := sq.DeleteTimeOutFromSqlite(time.Now().UnixNano())
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Info("周期性删除日志记录完成")
	})
	c.Start()
	return c
}
