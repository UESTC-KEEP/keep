package main

import (
	"keep/edge/cmd/edgeagent/app"
	"keep/edge/pkg/common/utils"
	kplogger "keep/pkg/util/kplogger"
	_ "net/http/pprof"
	"os"
)

func main() {
	defer utils.GracefulExit()

	kplogger.InitKeepLogger()

	command := app.NewEdgeAgentCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
