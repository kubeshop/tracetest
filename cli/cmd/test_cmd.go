package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Manage your tracetest tests",
	Long:  "Manage your tracetest tests",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Manage your tests")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
