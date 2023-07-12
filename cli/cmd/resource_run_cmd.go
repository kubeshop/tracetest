package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/runner"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

var (
	runParams = &runParameters{}
	runCmd    *cobra.Command
)

func init() {
	runCmd = &cobra.Command{
		GroupID: cmdGroupResources.ID,
		Use:     "run " + runnableResourceList(),
		Short:   "run resources",
		Long:    "run resources",
		PreRun:  setupCommand(),
		Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
			resourceType := resourceParams.ResourceName
			ctx := context.Background()

			r, err := runnerRegsitry.Get(resourceType)
			if err != nil {
				return "", fmt.Errorf("resource type '%s' cannot be run", resourceType)
			}

			orchestrator := runner.Orchestrator(
				cliLogger,
				utils.GetAPIClient(cliConfig),
				environmentClient,
			)

			runParams := runner.RunOptions{
				ID:              runParams.ID,
				DefinitionFile:  runParams.DefinitionFile,
				EnvID:           runParams.EnvID,
				SkipResultWait:  runParams.SkipResultWait,
				JUnitOuptutFile: runParams.JUnitOuptutFile,
			}

			exitCode, err := orchestrator.Run(ctx, r, runParams, output)
			if err != nil {
				return "", err
			}

			ExitCLI(exitCode)

			// ExitCLI will exit the process, so this return is just to satisfy the compiler
			return "", nil

		}, runParams),
		PostRun: teardownCommand,
	}

	runCmd.Flags().StringVarP(&runParams.DefinitionFile, "file", "f", "", "path to the definition file")
	runCmd.Flags().StringVar(&runParams.ID, "id", "", "id of the resource to run")
	runCmd.PersistentFlags().StringVarP(&runParams.EnvID, "environment", "e", "", "environment file or ID to be used")
	runCmd.PersistentFlags().BoolVarP(&runParams.SkipResultWait, "skip-result-wait", "W", false, "do not wait for results. exit immediately after test run started")
	runCmd.PersistentFlags().StringVarP(&runParams.JUnitOuptutFile, "junit", "j", "", "file path to save test results in junit format")
	rootCmd.AddCommand(runCmd)
}

type runParameters struct {
	ID              string
	DefinitionFile  string
	EnvID           string
	SkipResultWait  bool
	JUnitOuptutFile string
}

func (p runParameters) Validate(cmd *cobra.Command, args []string) []error {
	errs := []error{}
	if p.DefinitionFile == "" && p.ID == "" {
		errs = append(errs, paramError{
			Parameter: "resource",
			Message:   "you must specify a definition file or resource ID",
		})
	}

	if p.DefinitionFile != "" && p.ID != "" {
		errs = append(errs, paramError{
			Parameter: "resource",
			Message:   "you cannot specify both a definition file and resource ID",
		})
	}

	if p.JUnitOuptutFile != "" && p.SkipResultWait {
		errs = append(errs, paramError{
			Parameter: "junit",
			Message:   "--junit option is incompatible with --skip-result-wait option",
		})
	}

	return errs
}
