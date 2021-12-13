package filter

import (
	"fmt"
	beehiveContext "keep/pkg/util/core/context"
	//beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"keep/pkg/util/core/model"

	"keep/edge/pkg/common/modules"
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
			messsage.Content = log
			messsage.Router.Group = "/log"
			messsage.Router.Source = modules.LogsAgentModule
			//fmt.Println("+++++++++++++++++++++++  ", log)
			//if bufferpooler.PermissionOfSending {
			beehiveContext.Send(modules.EdgePublisherModule, *messsage)
			//}
		}
	case 7:
		if strings.Contains(log, "TRAC") {
			fmt.Println("OKKKKK")
		}
	}
}
