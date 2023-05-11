package misc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubeshop/tracetest/cli/global"
	misc_parameters "github.com/kubeshop/tracetest/cli/misc/parameters"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type DocGen struct {
	args[*misc_parameters.DocGen]
}

func NewDocGen(root global.Root) DocGen {
	defaults := NewDefaults("docgen", root.Setup)
	parameters := misc_parameters.NewDocGen()

	docGen := DocGen{
		args: NewArgs(defaults, parameters),
	}

	docGen.Cmd = &cobra.Command{
		GroupID: global.GroupMisc.ID,
		Use:     "docgen",
		Short:   "Generate the CLI documentation",
		Long:    "Generate the CLI documentation",
		Hidden:  true,
		PreRun:  defaults.PreRun,
		Run: defaults.Run(func(cmd *cobra.Command, args []string) (string, error) {
			os.RemoveAll(parameters.DocsOutputDir)
			err := os.MkdirAll(parameters.DocsOutputDir, os.ModePerm)
			if err != nil {
				fmt.Println(fmt.Errorf("could not create output dir: %w", err).Error())
				os.Exit(1)
			}

			err = doc.GenMarkdownTreeCustom(root.Cmd, parameters.DocsOutputDir, func(s string) string {
				return "# CLI Reference\n"
			}, func(s string) string { return s })
			if err != nil {
				fmt.Println(fmt.Errorf("could not generate documentation: %w", err).Error())
				os.Exit(1)
			}

			if parameters.DocusaurusFolder != "" {
				err = generateDocusaurusSidebar(parameters.DocsOutputDir, parameters.DocusaurusFolder)
				if err != nil {
					fmt.Println(fmt.Errorf("could not create docusaurus sidebar file: %w", err).Error())
					os.Exit(1)
				}
			}

			return "", nil
		}),
		PostRun: defaults.PostRun,
	}

	docGen.Cmd.Flags().StringVarP(&parameters.DocsOutputDir, "output", "o", "", "folder where all files will be generated")
	docGen.Cmd.Flags().StringVarP(&parameters.DocusaurusFolder, "docusaurus", "d", "", "folder containing your docusaurus documentation content")
	root.Cmd.AddCommand(docGen.Cmd)

	return docGen
}

func generateDocusaurusSidebar(outputDir string, docusaurusRootFolder string) error {
	fileContentTemplate := `
/** @type {import('@docusaurus/plugin-content-docs/lib/sidebars/types').SidebarItem[]} */
const pages = [
		%s
]

module.exports = pages;
`
	sidebarItemsContent, err := generateContentItems(outputDir, docusaurusRootFolder)
	if err != nil {
		return fmt.Errorf("could not list generated doc files: %w", err)
	}

	fileContent := fmt.Sprintf(fileContentTemplate, sidebarItemsContent)
	outputFile := filepath.Join(outputDir, "cli-sidebar.js")
	err = ioutil.WriteFile(outputFile, []byte(fileContent), 0644)
	if err != nil {
		return fmt.Errorf("could not write sidebar output file: %w", err)
	}

	return nil
}

func generateContentItems(inputDir string, docusaurusRootFolder string) (string, error) {
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return "", fmt.Errorf("could not read dir: %w", err)
	}

	entries := strings.Builder{}

	for _, file := range files {
		fileName := strings.TrimSuffix(file.Name(), ".md")
		command := strings.ReplaceAll(fileName, "_", " ")
		filePath, err := filepath.Rel(docusaurusRootFolder, filepath.Join(inputDir, fileName))
		if err != nil {
			return "", fmt.Errorf("could not get relative path: %w", err)
		}

		entry := fmt.Sprintf(`
		{
				type: "doc",
		label: "%s",
		id: "%s"
	},`, command, filePath)

		entries.Write([]byte(entry))
	}

	return entries.String(), nil
}
