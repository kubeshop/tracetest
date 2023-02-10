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

var (
	defaultSetters []func(*viper.Viper)
	flagSetters    = []func(*pflag.FlagSet){
		configFileFlag,
	}
)

func configFileFlag(flags *pflag.FlagSet) {
	flags.StringP("config", "c", "", "path to a config file")
}

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

func SetupFlags(flagset *pflag.FlagSet) {
	for _, fs := range flagSetters {
		fs(flagset)
	}
}

func New(flagset *pflag.FlagSet) (*Config, error) {
	vp := viper.New()
	vp.BindPFlags(flagset)

	configureConfigFile(vp)
	err := readConfigFile(vp)
	if err != nil {
		return nil, err
	}

	for _, ds := range defaultSetters {
		ds(vp)
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
