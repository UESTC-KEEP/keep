/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package logs

import (
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/spf13/cobra"
	"keep/constants"
	"keep/pkg/util/loggerv1.0.0"
	"time"
)

var logtopic string

// LogsCmd  represents the logs command
var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("logs called")
		switch logtopic {
		case "keep":
			printLogs(constants.DefaultEdgeLogFiles)
		case "cloudagent":
			printLogs(constants.DefaultCloudLogFiles)
		}
	},
}

func init() {
	LogsCmd.PersistentFlags().StringVarP(&logtopic, "agent", "a", "", "keep or cloudagent")
}

// 控制台打印日志文件
func printLogs(filename string) {
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
	for {
		line, ok = <-tails.Lines
		if !ok {
			logger.Error("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println(line.Text)
	}
}
