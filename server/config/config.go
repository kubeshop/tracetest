package config

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func SetupFlags(flags *pflag.FlagSet) {
	flags.StringP("config", "c", "", "path to a config file")
	configOptions.registerFlags(flags)
}

func New(confOpts ...Option) (*Config, error) {
	cfg := Config{
		vp: viper.New(),
	}

	for _, coFn := range confOpts {
		coFn(&cfg)
	}

	cfg.defaults()

	cfg.vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.vp.SetEnvPrefix("TRACETEST")
	cfg.vp.AutomaticEnv()

	cfg.configureConfigFile()

	configOptions.registerDefaults(cfg.vp)

	err := cfg.loadConfig()
	if err != nil {
		return nil, err
	}

	cfg.warnAboutDeprecatedFields()

	err = cfg.vp.Unmarshal(&cfg.config)
	if err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

type logger interface {
	Println(...any)
}

type Config struct {
	config    *oldConfig
	vp        *viper.Viper
	mu        sync.Mutex
	logger    logger
	resources resources
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

func (cfg *Config) loadConfig() error {
	if confFile := cfg.vp.GetString("config"); confFile != "" {
		// if --config is passed, and the file does not exists
		// it will trigger a "no such file or directory" error
		// which is NOT an instance of `viper.ConfigFileNotFoundError` checked later in this func
		// so this func WILL return an error
		cfg.vp.SetConfigFile(confFile)
	}

	err := cfg.vp.ReadInConfig()
	if path := cfg.vp.ConfigFileUsed(); path != "" {
		cfg.logger.Println("Config file used: ", path)
	}
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

func (cfg *Config) configureConfigFile() {
	cfg.vp.SetConfigName("tracetest")
	// intentionally removed this line, because it allows to have config files without extensions
	// cfg.vp.SetConfigType("yaml")
	cfg.vp.AddConfigPath("/etc/tracetest")
	cfg.vp.AddConfigPath("$HOME/.tracetest")
	cfg.vp.AddConfigPath(".")
}

func (c *Config) defaults() {
	if c.logger == nil {
		c.logger = log.Default()
	}
}

func (cfg *Config) warnAboutDeprecatedFields() error {
	for _, opt := range configOptions {
		if !opt.deprecated {
			continue
		}

		optionValue := cfg.vp.Get(opt.key)
		if optionValue == nil || optionValue == opt.defaultValue {
			continue
		}

		msg := fmt.Sprintf(`config "%s" is deprecated. `, opt.key)

		if opt.deprecationMessage != "" {
			msg += opt.deprecationMessage
		}

		cfg.logger.Println(msg)
	}

	return nil
}
