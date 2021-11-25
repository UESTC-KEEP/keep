package main

import (
	"keep/edge/cmd/edgeagent/app"
	"keep/edge/pkg/common/utils"
	"keep/pkg/util/kelogger"
	_ "net/http/pprof"
	"os"
)

func main() {
	defer utils.GracefulExit()

	kelogger.InitKeepLogger()

	command := app.NewEdgeAgentCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
