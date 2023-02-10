package cmd

import (
	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/spf13/cobra"
)

var (
	cfg         *config.Config
	appInstance *app.App

	rootCmd = &cobra.Command{
		Use:   "tracetest-server",
		Short: "tracetest server",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			appInstance, err = app.New(cfg)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	config.SetupFlags(rootCmd.PersistentFlags())

	cobra.OnInitialize(func() {
		var err error
		cfg, err = config.New(rootCmd.PersistentFlags())
		if err != nil {
			panic(err)
		}
	})

}
