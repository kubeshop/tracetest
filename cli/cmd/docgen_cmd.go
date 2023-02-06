package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsOutputDir string

var docGenCmd = &cobra.Command{
	Use:    "docgen",
	Short:  "Generate the CLI documentation",
	Long:   "Generate the CLI documentation",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		err := os.MkdirAll(docsOutputDir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		err = doc.GenMarkdownTreeCustom(rootCmd, docsOutputDir, func(s string) string {
			return "# CLI Reference\n"
		}, func(s string) string { return s })
		if err != nil {
			panic(err)
		}
	},
	PostRun: teardownCommand,
}

func init() {
	docGenCmd.PersistentFlags().StringVarP(&docsOutputDir, "output", "o", "", "tracetest docgen -o my/docs/dir")
	rootCmd.AddCommand(docGenCmd)
}
