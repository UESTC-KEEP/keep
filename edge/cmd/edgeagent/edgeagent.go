package main

import (
	"keep/edge/cmd/edgeagent/app"
	"keep/edge/pkg/common/utils"
	"keep/pkg/util/logs"
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
