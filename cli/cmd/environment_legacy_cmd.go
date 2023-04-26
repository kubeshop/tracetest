package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var environmentApplyFile string

var environmentCmd = &cobra.Command{
	GroupID: cmdGroupConfig.ID,
	Use:     "environment",
	Short:   "Manage your tracetest environments",
	Long:    "Manage your tracetest environments",
	PreRun:  setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PostRun: teardownCommand,
}

var environmentApplyCmd = &cobra.Command{
	Use:    "apply",
	Short:  "Create or update an environment to Tracetest",
	Long:   "Create or update an environment to Tracetest",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Warning! This is a deprecated command and it will be removed on Tracetest future versions!")
		fmt.Println("Please use `tracetest apply environment --file [path]` command instead.")
		fmt.Println("")

		// call new apply command
		definitionFile = dataStoreApplyFile
		applyCmd.Run(applyCmd, []string{"environment"})
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(environmentCmd)

	environmentApplyCmd.PersistentFlags().StringVarP(&environmentApplyFile, "file", "f", "", "file containing the environment configuration")
	environmentCmd.AddCommand(environmentApplyCmd)
}
