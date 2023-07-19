package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cliExitInterceptor func(code int)
)

var rootCmd = cobra.Command{
	Use:   "agent",
	Short: "Manages the tracetest agent",
	Long:  "Manages the tracetest agent",
}

func init() {
	rootCmd.AddCommand(&StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		ExitCLI(1)
	}
}

func ExitCLI(errorCode int) {
	if cliExitInterceptor != nil {
		cliExitInterceptor(errorCode)
		return
	}

	os.Exit(errorCode)
}

func RegisterCLIExitInterceptor(interceptor func(int)) {
	cliExitInterceptor = interceptor
}
