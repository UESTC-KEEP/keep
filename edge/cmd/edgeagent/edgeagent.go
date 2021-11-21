package main

import (
	"keep/edge/cmd/edgeagent/app"
<<<<<<< Updated upstream
	"keep/pkg/util/logs"
	_ "net/http/pprof"
=======
	"keep/edge/pkg/common/utils"
	"keep/pkg/util/kelogger"
>>>>>>> Stashed changes
	"os"
)

func main() {
<<<<<<< Updated upstream
	//defer utils.GracefulExit()
	logs.InitKeepLogger()
=======
	defer utils.GracefulExit()

	kelogger.InitKeepLogger()
>>>>>>> Stashed changes
	command := app.NewEdgeAgentCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
