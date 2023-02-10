package cmd

import (
	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/spf13/cobra"
)

var (
	cfg         *config.Config
	appInstance *app.App

	rootCmd = &cobra.Command{
		Use:   "tracetest-server",
		Short: "tracetest server",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			db, err := testdb.Connect(cfg.PostgresConnString())
			if err != nil {
				return err
			}

			appCfg := app.Config{
				Config: cfg,
			}

			appInstance, err = app.New(appCfg, db)
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
