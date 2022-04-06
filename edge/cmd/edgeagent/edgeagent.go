package main

import (
	"github.com/UESTC-KEEP/keep/edge/cmd/edgeagent/app"
	"github.com/UESTC-KEEP/keep/edge/pkg/common/utils"
	kplogger "github.com/UESTC-KEEP/keep/pkg/util/kplogger"
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
