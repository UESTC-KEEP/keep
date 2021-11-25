package main

import (
	"keep/cloud/cmd/cloudagent/app"
	"keep/cloud/pkg/common/utils"
	"keep/pkg/util/kelogger"
	"os"
)

func main() {
	command := app.NewCloudAgentCommand()
	kelogger.InitKeepLogger()
	defer utils.GracefulExit()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
