package cmd

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start tracete server",
	RunE: func(cmd *cobra.Command, args []string) error {
		var wg sync.WaitGroup
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			wg.Done()
			appInstance.Stop()
			os.Exit(1)
		}()

		wg.Add(1)
		err := appInstance.Start()
		if err != nil {
			return err
		}

		wg.Wait()

		return nil
	},
}
