package global_decorators

import (
	"fmt"
	"os"

	global_config "github.com/kubeshop/tracetest/cli/config"
	global_formatters "github.com/kubeshop/tracetest/cli/global/formatters"
	global_types "github.com/kubeshop/tracetest/cli/global/types"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type config struct {
	Logger
	Config     *global_config.Config
	Parameters *ConfigParams
}

type Config interface {
	Logger
	SetConfig(*global_config.Config)
	GetConfig() *global_config.Config
}

type ConfigParams struct {
	ConfigFile string
	Output     string

	// overrides
	OverrideEndpoint string
}

var _ Config = &config{}

func WithConfig(command global_types.Command) global_types.Command {
	logger, err := command.(Logger)
	if !err {
		panic("command must implement Logger interface")
	}

	config := &config{
		Logger: logger,
	}

	cmd := config.Get()
	cmd.PreRun = config.preRun(cmd.PreRun)

	cmd.PersistentFlags().StringVarP(&config.Parameters.Output, "output", "o", "", fmt.Sprintf("output format [%s]", global_formatters.OutputFormatsString))
	cmd.PersistentFlags().StringVarP(&config.Parameters.ConfigFile, "config", "c", "config.yml", "config file will be used by the CLI")
	cmd.PersistentFlags().StringVarP(&config.Parameters.OverrideEndpoint, "server-url", "s", "", "server url")
	config.Set(cmd)

	return config
}

func (d *config) SetConfig(cfg *global_config.Config) {
	d.Config = cfg
}

func (d *config) GetConfig() *global_config.Config {
	return d.Config
}

func (d *config) preRun(next CobraFn) CobraFn {
	return func(cmd *cobra.Command, args []string) {
		cfg, err := global_config.LoadConfig(d.Parameters.ConfigFile)
		if err != nil {
			d.GetLogger().Fatal("could not load config", zap.Error(err))
		}

		d.Config = &cfg

		if d.Parameters.OverrideEndpoint != "" {
			scheme, endpoint, err := global_config.ParseServerURL(d.Parameters.OverrideEndpoint)
			if err != nil {
				msg := fmt.Sprintf("cannot parse endpoint %s", d.Parameters.OverrideEndpoint)
				d.GetLogger().Error(msg, zap.Error(err))
				os.Exit(1)
			}

			d.Config.Scheme = scheme
			d.Config.Endpoint = endpoint
		}

		next(cmd, args)
	}
}
