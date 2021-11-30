package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"keep/constants"
	"keep/edge/cmd/edgeagent/app/options"
	"keep/edge/pkg/common/utils"
	"keep/edge/pkg/edgepublisher"
	"keep/edge/pkg/edgetwin"
	"keep/edge/pkg/healthzagent"
	"keep/edge/pkg/logsagent"
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	commonutil "keep/pkg/util"
	"keep/pkg/util/core"
	"keep/pkg/util/loggerv1.0.0"
	"net/http"
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
			// 性能监控
			go func() {
				logger.Debug(http.ListenAndServe(":6060", nil))
			}()
			config, err := opts.Config()
			text, err := yaml.Marshal(&config)
			// 写入配置文件
			err = ioutil.WriteFile(constants.DefaultEdgeagentConfigFile, text, 0777)
			// 下发配置文件
			commonutil.OrganizeConfigurationFile(constants.EdgeAgentName)
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}
			utils.PrintKEEPLogo()
			err = utils.EnvironmentCheck()
			//if err != nil {
			//	logger.Fatal(err)
			//	os.Exit(1)
			//}
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
	edgepublisher.Register(config.Modules.EdgePublisher, config.Modules.EdgePublisher.HostnameOverride, config.Modules.EdgePublisher.LocalIP)
	edgetwin.Register(config.Modules.EdgeTwin)
}
