package main

import (
	"keep/cloud/cmd/cloudagent/app"
	"keep/cloud/pkg/common/utils"
	"keep/pkg/util/logs"
	"os"
)

func main() {
	command := app.NewCloudAgentCommand()
	logs.InitKeepLogger()
	defer utils.GracefulExit()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
