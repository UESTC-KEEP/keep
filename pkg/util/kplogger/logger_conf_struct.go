package kplogger

import (
	"io"
	"time"
)

//配置结构体 复制粘贴的
type consoleLogger struct {
	Level    string `json:"level"`
	Colorful bool   `json:"color"`
	LogLevel int
}

type connLogger struct {
	innerWriter    io.WriteCloser
	ReconnectOnMsg bool   `json:"reconnectOnMsg"`
	Reconnect      bool   `json:"reconnect"`
	Net            string `json:"net"`
	Addr           string `json:"addr"`
	Level          string `json:"level"`
	LogLevel       int
	illNetFlag     bool //网络异常标记
}

type fileLogger struct {
	Filename   string `json:"filename"`
	Append     bool   `json:"append"`
	MaxLines   int    `json:"maxlines"`
	MaxSize    int    `json:"maxsize"`
	Daily      bool   `json:"daily"`
	MaxDays    int64  `json:"maxdays"`
	Level      string `json:"level"`
	PermitMask string `json:"permit"`

	LogLevel             int
	maxSizeCurSize       int
	maxLinesCurLines     int
	dailyOpenDate        int
	dailyOpenTime        time.Time
	fileNameOnly, suffix string
}

type logConfig struct {
	TimeFormat string         `json:"TimeFormat"`
	Console    *consoleLogger `json:"Console,omitempty"`
	File       *fileLogger    `json:"File,omitempty"`
	Conn       *connLogger    `json:"Conn,omitempty"`
}
