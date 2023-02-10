package config

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	config *config
	vp     *viper.Viper
	mu     sync.Mutex
}

type option struct {
	key          string
	defaultValue any
	description  string
}

type options []option

func (opts options) registerDefaults(vp *viper.Viper) {
	for _, opt := range opts {
		vp.SetDefault(opt.key, opt.defaultValue)
	}
}

func (opts options) registerFlags(flags *pflag.FlagSet) {
	for _, opt := range opts {
		switch defVal := opt.defaultValue.(type) {
		case int:
			flags.Int(opt.key, defVal, opt.description)
		case string:
			flags.String(opt.key, defVal, opt.description)
		case []string:
			flags.StringSlice(opt.key, defVal, opt.description)
		case bool:
			flags.Bool(opt.key, defVal, opt.description)
		default:
			panic(fmt.Errorf(
				"unexpected type %T in default value for config option %s",
				defVal, opt.key,
			))
		}

	}
}

var configOptions options

func configureConfigFile(vp *viper.Viper) {
	vp.SetConfigName("tracetest")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("/etc/tracetest")
	vp.AddConfigPath("$HOME/.tracetest")
	vp.AddConfigPath(".")
}

var ErrConfigFileNotFound = errors.New("config file not found")

func readConfigFile(vp *viper.Viper) error {
	if confFile := vp.GetString("config"); confFile != "" {
		// if --config is passed, and the file does not exists
		// it will trigger a "no such file or directory" error
		// which is NOT an instance of `viper.ConfigFileNotFoundError` checked later in this func
		// so this func WILL return an error
		vp.SetConfigFile(confFile)
	}

	err := vp.ReadInConfig()
	if err == nil {
		return nil
	}

	_, fileNotFound := err.(viper.ConfigFileNotFoundError)
	if fileNotFound {
		// config file is optional, can rely on defaults
		return nil
	}

	return fmt.Errorf("cannot read config file: %w", err)
}

func SetupFlags(flags *pflag.FlagSet) {
	flags.StringP("config", "c", "", "path to a config file")
	configOptions.registerFlags(flags)
}

func New(flags *pflag.FlagSet) (*Config, error) {
	vp := viper.New()

	configOptions.registerDefaults(vp)

	if flags != nil {
		vp.BindPFlags(flags)
	}

	configureConfigFile(vp)
	err := readConfigFile(vp)
	if err != nil {
		return nil, err
	}

	return &Config{
		config: &config{},
		vp:     vp,
	}, nil
}

func Must(c *Config, err error) *Config {
	if err != nil {
		panic(err)
	}

	return c
}

func (c *Config) AnalyticsEnabled() bool {
	if os.Getenv("TRACETEST_DEV") != "" {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.GA.Enabled
}
