package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information of Tracetest server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(appInstance.Version())
		fmt.Println("This is a temporary print")
	},
}
