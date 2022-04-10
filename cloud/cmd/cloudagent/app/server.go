package app

import (
	"github.com/UESTC-KEEP/keep/cloud/cmd/cloudagent/app/options"
	"github.com/UESTC-KEEP/keep/cloud/pkg/cloudimagemanager"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/informers"
	"github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient"
	"github.com/UESTC-KEEP/keep/cloud/pkg/promserver"
	"github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher"
	"github.com/UESTC-KEEP/keep/constants/cloud"
	"github.com/UESTC-KEEP/keep/edge/pkg/common/utils"
	cloudagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"github.com/UESTC-KEEP/keep/pkg/util/core"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"io/ioutil"
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
	//equalnodecontroller.Register(config.Modules.EqualNodeController)
	//tenantresourcequotacontroller.Register(config.Modules.TenantResourceQuotaController)
	//tenant_controller.Register(config.Modules.TenantController)
	promserver.Register(config.Modules.PromServer)
	requestDispatcher.Register(config.Modules.RequestDispatcher)
	cloudimageManager.Register(config.Modules.CloudImageManager)
}
