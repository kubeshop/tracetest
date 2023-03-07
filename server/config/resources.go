package config

import (
	"context"

	"github.com/kubeshop/tracetest/server/config/configresource"
)

type resources struct {
	config configResource
}

func WithConfigResource(cr configResource) Option {
	return func(c *Config) {
		c.resources.config = cr
	}
}

type configResource interface {
	Current(context.Context) configresource.Config
}
