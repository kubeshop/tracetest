package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	GroupID: cmdGroupCloud.ID,
	Use:     "login",
	Short:   "Tracetest Cloud Login",
	Long:    "Initializes the Tracetest Cloud Login process",
	PreRun:  setupCommand(),
	Run: WithResultHandler(func(_ *cobra.Command, _ []string) (string, error) {
		ctx := context.Background()
		onSuccess := func(token string, jwt string) {
			fmt.Printf("✔ Successful Authenticated with token: %s and jwt: %s\n", token, jwt)
			cliConfig.Token = token
			cliConfig.Jwt = jwt
			serverPath := ""
			cliConfig.ServerPath = &serverPath
			config.Save(ctx, cliConfig, config.ConfigureConfig{})
			ExitCLI(0)
		}

		onFailure := func(err error) {
			OnError(err)
		}

		// create server
		oauthServer = oauthServer.
			WithOnSuccess(onSuccess).
			WithOnFailure(onFailure)

		oauthServer.GetAuthCookie()
		return "", nil
	}),
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
