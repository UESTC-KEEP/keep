package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"keep/constants/device"
	"keep/constants/edge"
	"keep/device/cmd/deviceagent/app/options"
	devicemanager "keep/pkg/apis/compoenentconfig/keep/v1alpha1/devicemanager"
	commonutil "keep/pkg/util"
	"keep/pkg/util/core"
	logger "keep/pkg/util/loggerv1.0.1"
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

// NewDeviceCommand  create keep cmd
func NewDeviceCommand() *cobra.Command {
	opts := options.NewDefaultDeviceManagerConfig()
	cmd := &cobra.Command{
		Use:  "keep",
		Long: `keep description,however there is nothing in our code for now,so there is nothing in description`,
		Run: func(cmd *cobra.Command, args []string) {
			commonutil.OrganizeConfigurationFile(device.DeviceAgentName)
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
			commonutil.PrintKEEPLogo()
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
func registerModules(config *devicemanager.DeviceManagerConfig) {

}
