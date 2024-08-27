package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/telemetry"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&provisioningFile, "provisioning-file", "", "path to a provisioning file")
}

var provisioningFile string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Tracetest server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.Watch(func(updated *config.AppConfig) {
			appInstance.HotReload()
		})

		var wg sync.WaitGroup
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			wg.Done()
			appInstance.Stop()
			os.Exit(1)
		}()

		profilerConfig := cfg.ApplicationProfiler()
		if profilerConfig.Enabled {
			telemetry.StartProfiler(profilerConfig.Name, profilerConfig.Environment, profilerConfig.Endpoint, profilerConfig.SamplingRate)
		}

		wg.Add(1)
		err := appInstance.Start(app.WithProvisioningFile(provisioningFile))
		if err != nil {
			fmt.Println("Error starting server:", err.Error())
			return err
		}

		wg.Wait()

		return nil
	},
}
