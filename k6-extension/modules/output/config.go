package output

import (
	"fmt"
	"time"

	"go.k6.io/k6/output"
)

type Config struct {
	PushInterval time.Duration
}

func NewConfig(params output.Params) (Config, error) {
	cfg := Config{
		PushInterval: 1 * time.Second,
	}

	if val, ok := params.Environment["XK6_TRACETEST_PUSH_INTERVAL"]; ok {
		var err error
		cfg.PushInterval, err = time.ParseDuration(val)
		if err != nil {
			return cfg, fmt.Errorf("error parsing environment variable 'XK6_CROCOSPANS_PUSH_INTERVAL': %w", err)
		}
	}

	return cfg, nil
}
