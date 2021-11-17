package filter

import (
	"fmt"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/wonderivan/logger"
	beehiveContext "keep/core/context"
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgepublisher/bufferpooler"
	"keep/edge/pkg/logsagent/config"
	"strings"
	"time"
)

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
			messsage.Content = log
			bufferpooler.SendLogInQueue()
			go func() {
				resp, err := beehiveContext.SendSync(modules.EdgePublisherModule, *messsage, 5*time.Second)
				if err != nil {
					logger.Error(err)
				}
				fmt.Printf(modules.EdgePublisherModule+" 响应: %v, error: %v\n", resp, err)
				fmt.Println("发送日志至bufferpooler成功...")
			}()
		}
	case 7:
		if strings.Contains(log, "TRAC") {
			fmt.Println("OKKKKK")
		}
	}
}
