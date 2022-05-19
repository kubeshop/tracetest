package cmd

import (
	"log"
	"os"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var cliConfig config.Config
var cliLogger *zap.Logger

func setupCommand(cmd *cobra.Command, args []string) {
	loadConfig(cmd, args)
	setupLogger(cmd, args)
}

func loadConfig(cmd *cobra.Command, args []string) {
	config, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	cliConfig = config
}

func setupLogger(cmd *cobra.Command, args []string) {
	atom := zap.NewAtomicLevel()
	if verbose {
		atom.SetLevel(zap.DebugLevel)
	}

	encoderCfg := zap.NewProductionEncoderConfig()

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	cliLogger = logger
}

func teardownCommand(cmd *cobra.Command, args []string) {
	cliLogger.Sync()
}

func getAPIClient() *openapi.APIClient {
	config := openapi.NewConfiguration()
	config.Scheme = cliConfig.Scheme
	config.Host = cliConfig.Endpoint
	if cliConfig.ServerPath != nil {
		config.Servers = []openapi.ServerConfiguration{
			{
				URL: *cliConfig.ServerPath,
			},
		}
	}
	return openapi.NewAPIClient(config)
}
