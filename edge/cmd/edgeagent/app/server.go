package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"keep/constants"
	"keep/core"
	"keep/edge/cmd/edgeagent/app/options"
	"keep/edge/pkg/common/utils"
	"keep/edge/pkg/edgepublisher"
	"keep/edge/pkg/healthzagent"
	"keep/edge/pkg/logsagent"
<<<<<<< HEAD
	edgeagent "keep/pkg/apis/compoenentconfig/edgeagent/v1alpha1"
=======
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
>>>>>>> b0af266029c89d24fd39eac5960a66536ae9a802
	"os"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version called")
	},
}

// NewEdgeAgentCommand  create keep cmd
func NewEdgeAgentCommand() *cobra.Command {
	opts := options.NewDefaultEdgeAgentOptions()
	cmd := &cobra.Command{
		Use:  "keep",
		Long: `keep description,however there is nothing in our code for now,so there is nothing in description`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := opts.Config()
			text, err := yaml.Marshal(&config)
			// 写入配置文件
			err = ioutil.WriteFile(constants.DefaultEdgeagentConfigFile, text, 0777)
			if err != nil {
				logger.Fatal(err)
			}
			utils.PrintKEEPLogo()
			err = utils.EnvironmentCheck()
			if err != nil {
				logger.Fatal(err)
				os.Exit(1)
			}
			registerModules(config)
			core.Run()
		},
	}
	cmd.AddCommand(versionCmd)
	return cmd
}

// register all modules in system
func registerModules(config *edgeagent.EdgeAgentConfig) {
	healthzagent.Register(config.Modules.HealthzAgent)
	logsagent.Register(config.Modules.LogsAgent)
	edgepublisher.Register(config.Modules.EdgePublisher)
}
