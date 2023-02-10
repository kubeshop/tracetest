package config

import (
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	config *config
	vp     *viper.Viper
	mu     sync.Mutex
}

var defaultSetters []func(*viper.Viper)

func New() *Config {
	vp := viper.New()
	for _, ds := range defaultSetters {
		ds(vp)
	}
	return &Config{
		config: &config{},
		vp:     vp,
	}
}

func (c *Config) AnalyticsEnabled() bool {
	if os.Getenv("TRACETEST_DEV") != "" {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.GA.Enabled
}
