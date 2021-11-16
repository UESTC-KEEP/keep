package app

import (
<<<<<<< HEAD
	"github.com/wonderivan/logger"
	"github.com/spf13/cobra"
	"keep/cloud/cmd/cloudagent/app/options"
)

func NewCloudAgentCommnd() *cobra.Command {
	opts := options.NewDefaultCloudAgentOptions()
	cmd := &cobra.Command{
		Use: "cloudagent",
		Long: `cloudagent long description`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("cloudagent 开始启动！！！")
		}
	}
}
=======
	"github.com/spf13/cobra"
	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"keep/cloud/cmd/cloudagent/app/options"
	"keep/cloud/pkg/k8sclient"
	"keep/constants"
	"keep/core"
	"keep/edge/pkg/common/utils"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"os"
)

// NewCloudAgentCommand   create keep cmd
func NewCloudAgentCommand() *cobra.Command {
	opts := options.NewDefaultEdgeAgentOptions()
	cmd := &cobra.Command{
		Use:  "cloudagent",
		Long: `keep description,however there is nothing in our code for now,so there is nothing in description`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := opts.Config()
			text, err := yaml.Marshal(&config)
			// 写入配置文件
			err = ioutil.WriteFile(constants.DefaultCloudConfigFile, text, 0777)
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
	return cmd
}

// register all modules in system
func registerModules(config *cloudagent.CloudAgentConfig) {
	k8sclient.Register(config.Modules.K8sClient)
}
>>>>>>> bddbd7e0f200a771b61cbb6932118d2c7492d2c4
