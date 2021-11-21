package main

import (
	"keep/edge/cmd/edgeagent/app"
	"keep/pkg/util/logs"
	_ "net/http/pprof"
	"os"
)

func main() {
	//defer utils.GracefulExit()
	logs.InitKeepLogger()
	command := app.NewEdgeAgentCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
