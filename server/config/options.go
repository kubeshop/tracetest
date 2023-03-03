package config

import "github.com/spf13/pflag"

type Option func(*Config)

func WithFlagSet(flags *pflag.FlagSet) Option {
	return func(cfg *Config) {
		cfg.vp.BindPFlags(flags)
	}
}

func WithLogger(l logger) Option {
	return func(cfg *Config) {
		cfg.logger = l
	}
}
