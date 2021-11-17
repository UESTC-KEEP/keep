package main

import (
	"keep/edge/cmd/edgeagent/app"
	"keep/edge/pkg/common/utils"
	"keep/pkg/util/logs"
	"os"
)

func main() {
	defer utils.GracefulExit()

	logs.InitKeepLogger()
	command := app.NewEdgeAgentCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
