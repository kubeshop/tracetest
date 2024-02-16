package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsOutputDir string
var docusaurusFolder string

var docGenCmd = &cobra.Command{
	GroupID: cmdGroupMisc.ID,
	Use:     "docgen",
	Short:   "Generate the CLI documentation",
	Long:    "Generate the CLI documentation",
	Hidden:  true,
	PreRun:  setupLogger,
	Run: func(cmd *cobra.Command, args []string) {
		os.RemoveAll(docsOutputDir)
		err := os.MkdirAll(docsOutputDir, os.ModePerm)
		if err != nil {
			fmt.Println(fmt.Errorf("could not create output dir: %w", err).Error())
			ExitCLI(1)
		}

		err = doc.GenMarkdownTreeCustom(rootCmd, docsOutputDir, func(s string) string {
			return "# CLI Reference\n"
		}, func(s string) string { return s })
		if err != nil {
			fmt.Println(fmt.Errorf("could not generate documentation: %w", err).Error())
			ExitCLI(1)
		}

		if docusaurusFolder != "" {
			err = generateDocusaurusSidebar(docsOutputDir, docusaurusFolder)
			if err != nil {
				fmt.Println(fmt.Errorf("could not create docusaurus sidebar file: %w", err).Error())
				ExitCLI(1)
			}

			fixCodeBlocks(docsOutputDir)
		}
	},
	PostRun: teardownCommand,
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
	err = os.WriteFile(outputFile, []byte(fileContent), 0644)
	if err != nil {
		return fmt.Errorf("could not write sidebar output file: %w", err)
	}

	return nil
}

func generateContentItems(inputDir string, docusaurusRootFolder string) (string, error) {
	files, err := os.ReadDir(inputDir)
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

func fixCodeBlocks(inputDir string) error {
	path, err := filepath.Abs(inputDir)
	if err != nil {
		return fmt.Errorf("could not get absolute path: %w", err)
	}

	files, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("could not read dir: %w", err)
	}

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		if !strings.HasSuffix(filePath, ".md") {
			// ignore non-markdown files
			continue
		}

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("could not read file (%s) content: %w", filePath, err)
		}

		fileChanged := false
		lines := strings.Split(string(fileContent), "\n")
		for i, line := range lines {
			if len(line) > 0 && line[0] == '\t' {
				lines[i] = fmt.Sprintf("```\n%s\n```", strings.TrimPrefix(line, "\t"))
				fileChanged = true
			}
		}

		if fileChanged {
			fileContent := strings.Join(lines, "\n")
			err = os.WriteFile(filePath, []byte(fileContent), os.ModeAppend.Perm())
			if err != nil {
				return fmt.Errorf("could not update file %s: %w", filePath, err)
			}
		}
	}

	return nil
}

func init() {
	docGenCmd.Flags().StringVarP(&docsOutputDir, "output", "o", "", "folder where all files will be generated")
	docGenCmd.Flags().StringVarP(&docusaurusFolder, "docusaurus", "d", "", "folder containing your docusaurus documentation content")
	rootCmd.AddCommand(docGenCmd)
}
