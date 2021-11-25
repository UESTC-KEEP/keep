package filter

import (
	"fmt"
	beehiveContext "keep/pkg/util/core/context"
	"time"

	//beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"keep/pkg/util/core/model"

	"github.com/wonderivan/logger"
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgepublisher/bufferpooler"
	"keep/edge/pkg/logsagent/config"
	"strings"
)

var counter = 0

func FilterLogsByLevel(log string) {
	level := config.Config.LogLevel
	switch level {
	case 0:
		if strings.Contains(log, "EMER") {
			fmt.Println("OKKKKK")
		}
	case 1:
		if strings.Contains(log, "ALRT") {
			fmt.Println("OKKKKK")
		}
	case 2:
		if strings.Contains(log, "CRIT") {
			fmt.Println("OKKKKK")
		}
	case 3:
		if strings.Contains(log, "EROR") {
			fmt.Println("OKKKKK")
		}
	case 4:
		if strings.Contains(log, "WARN") {
			fmt.Println("OKKKKK")
		}
	case 5:
		if strings.Contains(log, "INFO") {
			fmt.Println("OKKKKK")
		}
	case 6:
		if strings.Contains(log, "DEBG") {
			//bufferpooler.SendLogInQueue(log)
			messsage := model.NewMessage("")
			counter++
			messsage.Content = log
			//fmt.Println("+++++++++++++++++++++++  ", log)
			if bufferpooler.PermissionOfSending {
				go func() {
					resp, err := beehiveContext.SendSync(modules.EdgePublisherModule, *messsage, 5*time.Second)
					if err != nil {
						logger.Error("logsagent消息发送超时:", err)
					} else {
						logger.Trace(modules.EdgePublisherModule+" 响应: %v, error: %v\n", resp, err)
						logger.Trace("发送日志至bufferpooler成功...")
					}
				}()
			}
		}
	case 7:
		if strings.Contains(log, "TRAC") {
			fmt.Println("OKKKKK")
		}
	}
}
