package cmd

import (
	"os"
	"strings"

	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	GroupID:    cmdGroupTests.ID,
	Use:        "test",
	Short:      "Manage your tracetest tests",
	Long:       "Manage your tracetest tests",
	Deprecated: "Please use `tracetest (apply|delete|export|get|list) test` commands instead.",
	PreRun:     setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PostRun: teardownCommand,
}

var testListCmd = &cobra.Command{
	Use:        "list",
	Short:      "List all tests",
	Long:       "List all tests",
	Deprecated: "Please use `tracetest list test` command instead.",
	PreRun:     setupCommand(),
	Run: func(_ *cobra.Command, _ []string) {
		listCmd.Run(listCmd, []string{"test"})
	},
	PostRun: teardownCommand,
}

var testExportCmd = &cobra.Command{
	Use:        "export",
	Short:      "Exports a test into a file",
	Long:       "Exports a test into a file",
	Deprecated: "Please use `tracetest export test` command instead.",
	PreRun:     setupCommand(),
	Run: func(_ *cobra.Command, _ []string) {
		exportCmd.Run(exportCmd, []string{"test"})
	},
	PostRun: teardownCommand,
}

var (
	runTestFileDefinition,
	runTestEnvID,
	runTestJUnit string
	runTestWaitForResult bool
)

var testRunCmd = &cobra.Command{
	Use:        "run",
	Short:      "Run a test on your Tracetest server",
	Long:       "Run a test on your Tracetest server",
	Deprecated: "Please use `tracetest run test` command instead.",
	PreRun:     setupCommand(),
	Run: func(_ *cobra.Command, _ []string) {
		// map old flags to new ones
		runParams.DefinitionFile = runTestFileDefinition
		runParams.EnvID = runTestEnvID
		runParams.SkipResultWait = !runTestWaitForResult
		runParams.JUnitOuptutFile = runTestJUnit

		fileType := getFileType(runTestFileDefinition)

		runCmd.Run(listCmd, []string{fileType})
	},
	PostRun: teardownCommand,
}

// getFileType returns the value of the `type` field defined in the definition file.
// if any error happens, it will return `test` as default, and errors will be ignored.
func getFileType(filePath string) string {
	const defaultFileType = "test"

	contents, err := os.ReadFile(filePath)
	if err != nil {
		return defaultFileType
	}

	var definition struct {
		Type string `json:"type"`
	}
	yaml := resourcemanager.Formats.Get(resourcemanager.FormatYAML)
	if err := yaml.Unmarshal(contents, &definition); err != nil {
		return defaultFileType
	}

	return strings.ToLower(definition.Type)
}

func init() {
	rootCmd.AddCommand(testCmd)

	// list
	testCmd.AddCommand(testListCmd)

	// export
	testExportCmd.PersistentFlags().StringVarP(&exportParams.ResourceID, "id", "", "", "id of the test")
	testExportCmd.PersistentFlags().StringVarP(&exportParams.OutputFile, "output", "o", "", "file to be created with definition")
	testCmd.AddCommand(testExportCmd)

	// run
	testRunCmd.PersistentFlags().StringVarP(&runTestEnvID, "environment", "e", "", "id of the environment to be used")
	testRunCmd.PersistentFlags().StringVarP(&runTestFileDefinition, "definition", "d", "", "path to the definition file to be run")
	testRunCmd.PersistentFlags().BoolVarP(&runTestWaitForResult, "wait-for-result", "w", false, "wait for the test result to print it's result")
	testRunCmd.PersistentFlags().StringVarP(&runTestJUnit, "junit", "j", "", "path to the junit file that will be generated")

	testCmd.AddCommand(testRunCmd)
}
