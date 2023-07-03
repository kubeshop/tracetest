package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	GroupID: cmdGroupMisc.ID,
	Use:     "dashboard",
	Short:   "Opens the Tracetest Dashboard URL",
	Long:    "Opens the Tracetest Dashboard URL",
	PreRun:  setupCommand(),
	Run: WithResultHandler(func(_ *cobra.Command, _ []string) (string, error) {
		if cliConfig.IsEmpty() {
			return "", fmt.Errorf("missing Tracetest endpoint configuration")
		}

		err := openBrowser(cliConfig.URL())
		if err != nil {
			return "", fmt.Errorf("failed to open the dashboard url: %s", cliConfig.URL())
		}

		return fmt.Sprintf("Opening \"%s\" in the default browser", cliConfig.URL()), nil
	}),
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}

func openBrowser(u string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", u).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start()
	case "darwin":
		return exec.Command("open", u).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
