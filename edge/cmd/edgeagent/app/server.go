package app

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-ps"
	"github.com/spf13/cobra"
	"github.com/wonderivan/logger"
	"keep/core"
	"keep/edge/cmd/edgeagent/app/options"
	"keep/edge/pkg/healthzagent"
	edgeagent "keep/pkg/apis/compoenentconfig/edgeagent/v1alpha1"
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

// NewEdgeAgentCommand  create edgeagent cmd
func NewEdgeAgentCommand() *cobra.Command {
	opts := options.NewDefaultEdgeAgentOptions()
	cmd := &cobra.Command{
		Use:  "edgeagent",
		Long: `edgeagent description,however there is nothing in our code for now,so there is nothing in description`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := opts.Config()
			if err != nil {
				logger.Fatal(err)
			}
			registerModules(config)
			logger.Info("命令创建成功")
			core.Run()
		},
	}
	cmd.AddCommand(versionCmd)
	return cmd
}

// environmentCheck check the environment before edgeagent start
// if Check failed,  return errors
func environmentCheck() error {
	// if kubelet is running, return error
	if find, err := findProcess("kubelet"); err != nil {
		return err
	} else if find {
		return errors.New("Kubelet should not running on edge node when running edgeagent")
	}

	// if kube-proxy is running, return error
	if find, err := findProcess("kube-proxy"); err != nil {
		return err
	} else if find {
		return errors.New("Kube-proxy should not running on edge node when running edgeagent")
	}

	return nil
}

// findProcess find a running process by name
func findProcess(name string) (bool, error) {
	processes, err := ps.Processes()
	if err != nil {
		return false, err
	}

	for _, process := range processes {
		if process.Executable() == name {
			return true, nil
		}
	}

	return false, nil
}

// register all modules in system
func registerModules(config *edgeagent.EdgeAgentConfig) {
	healthzagent.Register(config.Modules.HealthzAgent)
}
