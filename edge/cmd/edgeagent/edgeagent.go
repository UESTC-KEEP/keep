package main

import (
	"keep/edge/cmd/edgeagent/app"
	"keep/edge/pkg/common/logs"
	"keep/edge/pkg/common/utils"
	"os"
)

func main() {
	command := app.NewEdgeAgentCommand()
	logs.InitKeepLogger()
	defer utils.GracefulExit()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
