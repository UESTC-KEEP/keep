package app

import (
	"flag"
	"keep/kpadm/cmd/kpadm/app/cmd"
	"os"

	"github.com/spf13/pflag"
)

//Run executes the keadm command
func Run() error {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	cmd := cmd.NewKeepEdgeCommand(os.Stdin, os.Stdout, os.Stderr)
	return cmd.Execute()
}
