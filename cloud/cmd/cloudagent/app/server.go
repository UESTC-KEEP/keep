package app

import (
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