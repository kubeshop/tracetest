package config

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
)

func validateDuration(key string) func(c *AppConfig) error {
	return func(c *AppConfig) error {
		input := c.vp.GetString(key)
		_, err := time.ParseDuration(input)
		if err == nil {
			return nil
		}

		return fmt.Errorf("invalid duration format '%s'", input)
	}
}

func (c *AppConfig) Validate() error {
	var err error
	for _, opt := range configOptions {
		if opt.validate == nil {
			// no validator defined
			continue
		}

		optErr := opt.validate(c)
		if optErr == nil {
			// no error, move on
			continue
		}

		err = multierror.Append(
			err,
			fmt.Errorf("invalid config value for '%s': %s", opt.key, optErr.Error()),
		)
	}

	return err
}
