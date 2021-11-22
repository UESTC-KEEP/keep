package main

import (
	"keep/edge/cmd/edgeagent/app"
	"keep/pkg/util/logs"
	_ "net/http/pprof"
	"os"
)

func main() {
	//defer utils.GracefulExit()
	command := app.NewEdgeAgentCommand()
	logs.InitKeepLogger()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
