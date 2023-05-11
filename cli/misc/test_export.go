package misc

import (
	"context"

	"github.com/kubeshop/tracetest/cli/analytics"
	misc_actions "github.com/kubeshop/tracetest/cli/misc/actions"
	misc_parameters "github.com/kubeshop/tracetest/cli/misc/parameters"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type TestExport struct {
	args[*misc_parameters.TestExport]
}

func (t TestExport) Run(cmd *cobra.Command, args []string) (string, error) {
	analytics.Track("Test Export", "cmd", map[string]string{})

	ctx := context.Background()
	t.Setup.Logger.Debug("Exporting test", zap.String("testID", t.Parameters.ExportTestId))
	exportTestAction := misc_actions.NewExportTestAction(*t.Setup.Config, t.Setup.Logger, t.Setup.Client)

	actionArgs := misc_actions.ExportTestConfig{
		TestId:     t.Parameters.ExportTestId,
		OutputFile: t.Parameters.ExportTestOutputFile,
		Version:    t.Parameters.Version,
	}

	err := exportTestAction.Run(ctx, actionArgs)
	return "", err
}

func NewTestExport(root Test) TestExport {
	parameters := misc_parameters.NewTestExport()
	defaults := NewDefaults("test export", root.Setup.Setup)

	testExport := TestExport{
		args: NewArgs(defaults, parameters),
	}

	testExport.Cmd = &cobra.Command{
		Use:     "export",
		Short:   "Exports a test into a file",
		Long:    "Exports a test into a file",
		PreRun:  defaults.PreRun,
		Run:     defaults.Run(testExport.Run),
		PostRun: defaults.PostRun,
	}

	testExport.Cmd.PersistentFlags().StringVarP(&parameters.ExportTestId, "id", "", "", "id of the test")
	testExport.Cmd.PersistentFlags().StringVarP(&parameters.ExportTestOutputFile, "output", "o", "", "file to be created with definition")
	testExport.Cmd.PersistentFlags().Int32VarP(&parameters.Version, "version", "", -1, "version of the test. Default is latest")

	root.Cmd.AddCommand(testExport.Cmd)

	return testExport
}
