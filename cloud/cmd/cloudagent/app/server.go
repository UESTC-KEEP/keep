package app

import (
	"io/ioutil"
	"keep/cloud/cmd/cloudagent/app/options"
	"keep/cloud/pkg/cloudimagemanager"
	"keep/cloud/pkg/common/client"
	"keep/cloud/pkg/common/informers"
	"keep/cloud/pkg/equalnodecontroller"
	"keep/cloud/pkg/k8sclient"
	"keep/cloud/pkg/promserver"
	"keep/cloud/pkg/requestDispatcher"
	"keep/cloud/pkg/tenantresourcequotacontroller"
	"keep/constants/cloud"
	"keep/edge/pkg/common/utils"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"keep/pkg/util/core"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/loggerv1.0.1"
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
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
			config, err := opts.Config()
			text, err := yaml.Marshal(&config)
			// 写入配置文件
			err = ioutil.WriteFile(cloud.DefaultCloudConfigFile, text, 0777)
			if err != nil {
				logger.Fatal(err)
			}
			// 初始化口k8s配置
			client.InitKubeEdgeClient(config.Modules.K8sClient)

			gis := informers.GetInformersManager()

			utils.PrintKEEPLogo()
			// err = utils.EnvironmentCheck()
			// if err != nil {
			// 	logger.Fatal(err)
			// 	os.Exit(1)
			// }
			registerModules(config)
			core.StartModules()
			gis.Start(beehiveContext.Done())
			core.GracefulShutdown()
		},
	}
	return cmd
}

// register all modules in system
func registerModules(config *cloudagent.CloudAgentConfig) {
	k8sclient.Register(config.Modules.K8sClient)
	equalnodecontroller.Register(config.Modules.EqualNodeController)
	tenantresourcequotacontroller.Register(config.Modules.TenantResourceQuotaController)
	promserver.Register(config.Modules.PromServer)
	requestDispatcher.Register(config.Modules.RequestDispatcher)
	cloudimageManager.Register(config.Modules.CloudImageManager)
}
