package tailf

import (
	"github.com/UESTC-KEEP/keep/edge/pkg/logsagent/filter"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"github.com/hpcloud/tail"
	"time"
)

func watchLogFileLine(filename string) {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,                                 // 监听新行，使用tail -f，这个参数非常重要
	}
	tails, err := tail.TailFile(filename, config)
	if err != nil {
		logger.Fatal("监听文件 "+filename+" 失败,退出监听 err:", err)
		return
	}
	var line *tail.Line
	var ok bool
	logger.Info("启动监听文件:", filename)
	for {
		line, ok = <-tails.Lines
		if !ok {
			logger.Error("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		filter.FilterLogsByLevel(line.Text)
	}
}

// StartWatchingLogs 对配置中定义的需要进行监听的文件 分别启动携程
func StartWatchingLogs(logfiles []string) {
	for i := 0; i < len(logfiles); i++ {
		go watchLogFileLine(logfiles[i])
	}
}
