package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/cloud/cmd"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/spf13/cobra"
)

var waitParams = &waitParameters{}

var waitCmd = &cobra.Command{
	GroupID: cmdGroupMisc.ID,
	Use:     "wait group",
	Short:   "Waits for a run group to be finished",
	Long:    "Waits for a run group to be finished and displays the result",
	PreRun:  setupCommand(),
	Run: WithResultHandler(WithParamsHandler(waitParams)(func(_ context.Context, _ *cobra.Command, _ []string) (string, error) {
		exitCode, err := cmd.Wait(context.Background(), &cliConfig, waitParams.RunGroupID, output)

		ExitCLI(exitCode)
		return "", err
	})),
	PostRun: teardownCommand,
}

func init() {
	waitCmd.Flags().StringVarP(&waitParams.RunGroupID, "id", "", "", "Run Group ID to wait")

	rootCmd.AddCommand(waitCmd)
}

type waitParameters struct {
	RunGroupID string
}

func (p waitParameters) Validate(cmd *cobra.Command, args []string) []error {
	errors := make([]error, 0)

	if p.RunGroupID == "" {
		errors = append(errors, cmdutil.ParamError{
			Parameter: "group",
			Message:   "Run Group ID is required",
		})
	}

	return errors
}
