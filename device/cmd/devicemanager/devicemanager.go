package main

import (
	"keep/device/cmd/devicemanager/app"
	"keep/pkg/util/kplogger"
	"os"
)

func main() {
	kplogger.InitKeepLogger()
	command := app.NewDeviceCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
