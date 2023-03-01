package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type (
	googleAnalytics struct {
		Enabled bool `yaml:",omitempty" mapstructure:"enabled"`
	}

	oldConfig struct {
		Server    serverConfig    `yaml:",omitempty" mapstructure:"server"`
		GA        googleAnalytics `yaml:"googleAnalytics,omitempty" mapstructure:"googleAnalytics"`
		Telemetry telemetry       `yaml:",omitempty" mapstructure:"telemetry"`
		Demo      demo            `yaml:",omitempty" mapstructure:"demo"`
	}
)

type Config struct {
	config *oldConfig
	vp     *viper.Viper
	mu     sync.Mutex
}

type logger interface {
	Println(...any)
}

type option struct {
	key                string
	defaultValue       any
	description        string
	validate           func(*Config) error
	deprecated         bool
	deprecationMessage string
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
	if fileNotFound || vp.GetString("config") == "" {
		// config file is optional, can rely on defaults
		return nil
	}

	return fmt.Errorf("cannot read config file: %w", err)
}

func SetupFlags(flags *pflag.FlagSet) {
	flags.StringP("config", "c", "", "path to a config file")
	configOptions.registerFlags(flags)
}

func New(flags *pflag.FlagSet, logger logger) (*Config, error) {
	vp := viper.New()

	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vp.SetEnvPrefix("TRACETEST")
	vp.AutomaticEnv()

	configureConfigFile(vp)

	if flags != nil {
		vp.BindPFlags(flags)
	}

	configOptions.registerDefaults(vp)

	warnAboutDeprecatedFields(vp, logger)

	err := readConfigFile(vp)
	if err != nil {
		return nil, err
	}

	oldConfig := oldConfig{}
	err = vp.Unmarshal(&oldConfig)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		config: &oldConfig,
		vp:     vp,
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func Must(c *Config, err error) *Config {
	if err != nil {
		panic(err)
	}

	return c
}

func (c *Config) Watch(updateFn func(c *Config)) {
	c.vp.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		updateFn(c)
	})
	c.vp.WatchConfig()
}

func (c *Config) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.vp.Set(key, value)
}

func (c *Config) AnalyticsEnabled() bool {
	if os.Getenv("TRACETEST_DEV") != "" {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.GA.Enabled
}

func warnAboutDeprecatedFields(vp *viper.Viper, logger logger) error {
	err := readConfigFile(vp)
	if err != nil {
		return err
	}

	for _, opt := range configOptions {
		if !opt.deprecated {
			continue
		}

		optionValue := vp.Get(opt.key)
		if optionValue == nil || optionValue == opt.defaultValue {
			continue
		}

		msg := fmt.Sprintf(`config "%s" is deprecated. `, opt.key)

		if opt.deprecationMessage != "" {
			msg += opt.deprecationMessage
		}

		logger.Println(msg)
	}

	return nil
}
