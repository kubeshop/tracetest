package pollingprofile

import (
	"fmt"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
)

type Strategy string

const (
	Periodic Strategy = "periodic"
)

const (
	ResourceName       = "PollingProfile"
	ResourceNamePlural = "PollingProfiles"
)

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationGet,
	resourcemanager.OperationList,
	resourcemanager.OperationUpdate,
}

var DefaultPollingProfile = PollingProfile{
	ID:       id.ID("current"),
	Name:     "default",
	Default:  true,
	Strategy: Periodic,
	Periodic: &PeriodicPollingConfig{
		Timeout:    "1m",
		RetryDelay: "5s",
	},
}

type PollingProfile struct {
	ID       id.ID                  `json:"id"`
	Name     string                 `json:"name"`
	Default  bool                   `json:"default"`
	Strategy Strategy               `json:"strategy"`
	Periodic *PeriodicPollingConfig `json:"periodic"`
}

type PeriodicPollingConfig struct {
	RetryDelay           string `json:"retryDelay"`
	Timeout              string `json:"timeout"`
	SelectorMatchRetries int    `json:"selectorMatchRetries"`
}

func (ppc *PeriodicPollingConfig) TimeoutDuration() time.Duration {
	d, _ := time.ParseDuration(ppc.Timeout)
	return d
}

func (ppc *PeriodicPollingConfig) RetryDelayDuration() time.Duration {
	d, _ := time.ParseDuration(ppc.RetryDelay)
	return d
}

func (ppc *PeriodicPollingConfig) MaxTracePollRetry() int {
	return int(math.Ceil(float64(ppc.TimeoutDuration()) / float64(ppc.RetryDelayDuration())))
}

func (ppc *PeriodicPollingConfig) Validate() error {
	if ppc == nil {
		return fmt.Errorf("missing periodic polling profile configuration")
	}

	if _, err := time.ParseDuration(ppc.RetryDelay); err != nil {
		return fmt.Errorf("retry delay configuration is invalid: %w", err)
	}

	if _, err := time.ParseDuration(ppc.Timeout); err != nil {
		return fmt.Errorf("timeout configuration is invalid: %w", err)
	}

	return nil
}

func (pp PollingProfile) HasID() bool {
	return pp.ID.String() != ""
}

func (pp PollingProfile) GetID() id.ID {
	return pp.ID
}

func (pp PollingProfile) Validate() error {
	if pp.Strategy == Periodic {
		if err := pp.Periodic.Validate(); err != nil {
			return err
		}
	}

	return nil
}
