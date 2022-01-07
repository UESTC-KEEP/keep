package main

import (
	"github.com/UESTC-KEEP/keep/cloud/cmd/cloudagent/app"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/utils"
	"github.com/UESTC-KEEP/keep/constants/cloud"
	commonutil "github.com/UESTC-KEEP/keep/pkg/util"
	"github.com/UESTC-KEEP/keep/pkg/util/kplogger"
	"os"
)

func main() {

	command := app.NewCloudAgentCommand()
	commonutil.OrganizeConfigurationFile(cloud.CloudAgentName)
	kplogger.InitKeepLogger()
	defer utils.GracefulExit()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
