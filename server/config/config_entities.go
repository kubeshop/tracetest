package config

import (
	"os"

	"github.com/kubeshop/tracetest/server/pkg/id"
)

type Config struct {
	ID   id.ID  `json:"id"`
	Name string `json:"name"`

	AnalyticsEnabled bool `json:"analyticsEnabled"`
}

func (c Config) HasID() bool {
	return c.ID.String() != ""
}

func (c Config) GetID() id.ID {
	return c.ID
}

func (c Config) Validate() error {
	return nil
}

func (c Config) IsAnalyticsEnabled() bool {
	if os.Getenv("TRACETEST_DEV") != "" {
		return false
	}

	return c.AnalyticsEnabled
}
