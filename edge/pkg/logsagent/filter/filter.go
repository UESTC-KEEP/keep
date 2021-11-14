package filter

import (
	"fmt"
	"keep/edge/pkg/logsagent/config"
	"strings"
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
			fmt.Println("OKKKKK")
		}
	case 7:
		if strings.Contains(log, "TRAC") {
			fmt.Println("OKKKKK")
		}
	}
}
