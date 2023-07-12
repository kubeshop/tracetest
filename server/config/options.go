package config

import "github.com/spf13/pflag"

type Option func(*AppConfig)

func WithFlagSet(flags *pflag.FlagSet) Option {
	return func(cfg *AppConfig) {
		cfg.vp.BindPFlags(flags)
	}
}

func WithLogger(l logger) Option {
	return func(cfg *AppConfig) {
		cfg.logger = l
	}
}
