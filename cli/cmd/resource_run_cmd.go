package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/runner"
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
		PreRun:  setupCommand(WithOptionalResourceName()),
		Run: WithResourceMiddleware(func(_ *cobra.Command, args []string) (string, error) {
			ctx := context.Background()
			resourceType, err := getResourceType(runParams, resourceParams)
			if err != nil {
				return "", err
			}

			r, err := runnerRegistry.Get(resourceType)
			if err != nil {
				return "", fmt.Errorf("resource type '%s' cannot be run", resourceType)
			}

			orchestrator := runner.Orchestrator(
				cliLogger,
				config.GetAPIClient(cliConfig),
				variableSetClient,
			)

			if runParams.EnvID != "" {
				runParams.VarsID = runParams.EnvID
			}

			runParams := runner.RunOptions{
				ID:              runParams.ID,
				DefinitionFile:  runParams.DefinitionFile,
				VarsID:          runParams.VarsID,
				SkipResultWait:  runParams.SkipResultWait,
				JUnitOuptutFile: runParams.JUnitOuptutFile,
				RequiredGates:   runParams.RequriedGates,
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
	runCmd.Flags().StringVarP(&runParams.VarsID, "vars", "", "", "variable set file or ID to be used")
	runCmd.Flags().BoolVarP(&runParams.SkipResultWait, "skip-result-wait", "W", false, "do not wait for results. exit immediately after test run started")
	runCmd.Flags().StringVarP(&runParams.JUnitOuptutFile, "junit", "j", "", "file path to save test results in junit format")
	runCmd.Flags().StringSliceVar(&runParams.RequriedGates, "required-gates", []string{}, "override default required gate. "+validRequiredGatesMsg())

	//deprecated
	runCmd.Flags().StringVarP(&runParams.EnvID, "environment", "e", "", "environment file or ID to be used")
	runCmd.Flags().MarkDeprecated("environment", "use --vars instead")
	runCmd.Flags().MarkShorthandDeprecated("e", "use --vars instead")

	rootCmd.AddCommand(runCmd)
}

func getResourceType(runParams *runParameters, resourceParams *resourceParameters) (string, error) {
	if resourceParams.ResourceName != "" {
		return resourceParams.ResourceName, nil
	}

	if runParams.DefinitionFile != "" {
		filePath := runParams.DefinitionFile
		f, err := fileutil.Read(filePath)
		if err != nil {
			return "", fmt.Errorf("cannot read file %s: %w", filePath, err)
		}

		return strings.ToLower(f.Type()), nil
	}

	return "", fmt.Errorf("resourceName is empty and no definition file provided")

}

func validRequiredGatesMsg() string {
	opts := make([]string, 0, len(openapi.AllowedSupportedGatesEnumValues))
	for _, v := range openapi.AllowedSupportedGatesEnumValues {
		opts = append(opts, string(v))
	}

	return "valid options: " + strings.Join(opts, ", ")
}

type runParameters struct {
	ID              string
	DefinitionFile  string
	VarsID          string
	EnvID           string
	SkipResultWait  bool
	JUnitOuptutFile string
	RequriedGates   []string
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

	for _, rg := range p.RequriedGates {
		_, err := openapi.NewSupportedGatesFromValue(rg)
		if err != nil {
			errs = append(errs, paramError{
				Parameter: "required-gates",
				Message:   fmt.Sprintf("invalid option '%s'. "+validRequiredGatesMsg(), rg),
			})
		}
	}

	return errs
}
