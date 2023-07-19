package cmd

import (
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/initialization"
	"github.com/spf13/cobra"
)

var (
	apiKey  string
	devMode bool
)

var StartCmd = cobra.Command{
	Use:   "start",
	Short: "Start the local agent",
	Long:  "Start the local agent",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			ExitCLI(1)
		}

		initialization.Start(cfg)
	},
}

func init() {
	StartCmd.Flags().StringVarP(&apiKey, "apiKey", "", "", "the API key from the environment that will run the tests")
	StartCmd.Flags().BoolVarP(&devMode, "devMode", "d", false, "starts a dev mode session on your private environment")
}
