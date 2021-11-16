package main

import (
	"keep/cloud/cmd/cloudagent/app"
<<<<<<< HEAD
	"keep/cloud/pkg/common/logs"
	"keep/edge/pkg/common/utils"
=======
	"keep/edge/pkg/common/utils"
	"keep/pkg/util/logs"
>>>>>>> bddbd7e0f200a771b61cbb6932118d2c7492d2c4
	"os"
)

func main() {
<<<<<<< HEAD
	command := app.NewCloudAgentCommnd()
=======
	command := app.NewCloudAgentCommand()
>>>>>>> bddbd7e0f200a771b61cbb6932118d2c7492d2c4
	logs.InitKeepLogger()
	defer utils.GracefulExit()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
<<<<<<< HEAD
=======

>>>>>>> bddbd7e0f200a771b61cbb6932118d2c7492d2c4
