package main

import (
	"keep/cloud/cmd/cloudagent/app"
	"keep/cloud/pkg/common/logs"
	"keep/edge/pkg/common/utils"
	"os"
)

func main() {
	command := app.NewCloudAgentCommnd()
	logs.InitKeepLogger()
	defer utils.GracefulExit()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
