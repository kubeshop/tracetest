package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Must(c *AppConfig, err error) *AppConfig {
	if err != nil {
		panic(err)
	}

	return c
}

type appConfigOption struct {
	key                string
	defaultValue       any
	description        string
	validate           func(*AppConfig) error
	deprecated         bool
	deprecationMessage string
}

type options []appConfigOption

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
