package app

import (
	"io/ioutil"
	"keep/cloud/cmd/cloudagent/app/options"
	"keep/cloud/pkg/k8sclient"
	"keep/cloud/pkg/promserver"
	"keep/constants"
	"keep/edge/pkg/common/utils"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	commonutil "keep/pkg/util"
	"keep/pkg/util/core"
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v2"
)

// NewCloudAgentCommand   create keep cmd
func NewCloudAgentCommand() *cobra.Command {
	opts := options.NewDefaultEdgeAgentOptions()
	cmd := &cobra.Command{
		Use:  "cloudagent",
		Long: `keep description,however there is nothing in our code for now,so there is nothing in description`,
		Run: func(cmd *cobra.Command, args []string) {
			// 性能监控
			go func() {
				logger.Debug(http.ListenAndServe(":6060", nil))
			}()
			commonutil.OrganizeConfigurationFile(constants.CloudAgentName)
			config, err := opts.Config()
			text, err := yaml.Marshal(&config)
			// 写入配置文件
			err = ioutil.WriteFile(constants.DefaultCloudConfigFile, text, 0777)
			if err != nil {
				logger.Fatal(err)
			}
			utils.PrintKEEPLogo()
			// err = utils.EnvironmentCheck()
			// if err != nil {
			// 	logger.Fatal(err)
			// 	os.Exit(1)
			// }
			registerModules(config)
			core.Run()
		},
	}
	return cmd
}

// register all modules in system
func registerModules(config *cloudagent.CloudAgentConfig) {
	k8sclient.Register(config.Modules.K8sClient)
	promserver.Register(config.Modules.PromServer)
}
