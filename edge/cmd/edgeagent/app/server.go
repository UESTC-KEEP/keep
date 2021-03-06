package app

import (
	"fmt"
	"github.com/UESTC-KEEP/keep/constants/edge"
	"github.com/UESTC-KEEP/keep/edge/cmd/edgeagent/app/options"
	"github.com/UESTC-KEEP/keep/edge/pkg/common/utils"
	dmi "github.com/UESTC-KEEP/keep/edge/pkg/device_manage_interface"
	"github.com/UESTC-KEEP/keep/edge/pkg/edgepublisher"
	"github.com/UESTC-KEEP/keep/edge/pkg/edgetwin"
	"github.com/UESTC-KEEP/keep/edge/pkg/healthzagent"
	"github.com/UESTC-KEEP/keep/edge/pkg/logsagent"
	edgeagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	commonutil "github.com/UESTC-KEEP/keep/pkg/util"
	"github.com/UESTC-KEEP/keep/pkg/util/core"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
			//go func() {
			//	logger.Debug(http.ListenAndServe(":6060", nil))
			//}()
			config, err := opts.Config()
			text, err := yaml.Marshal(&config)
			// 写入配置文件
			err = ioutil.WriteFile(edge.DefaultEdgeagentConfigFile, text, 0777)
			// 下发配置文件
			commonutil.OrganizeConfigurationFile(edge.EdgeAgentName)
			if err != nil {
				logger.Error(err)
				os.Exit(1)
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
	edgetwin.Register(config.Modules.EdgeTwin)
	dmi.Register(config.Modules.DeviceMapperInterface)
}
