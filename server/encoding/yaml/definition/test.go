package definition

import (
	"fmt"
)

type Test struct {
	ID          string      `mapstructure:"id"`
	Name        string      `mapstructure:"name"`
	Description string      `mapstructure:"description"`
	Trigger     TestTrigger `mapstructure:"trigger"`
	Specs       []TestSpec  `mapstructure:"specs,omitempty"`
	Outputs     []Output    `mapstructure:"outputs,omitempty"`
}

func (t Test) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("test name cannot be empty")
	}

	if err := t.Trigger.Validate(); err != nil {
		return fmt.Errorf("test trigger must be valid: %w", err)
	}

	return nil
}

type TestTrigger struct {
	Type        string      `mapstructure:"type"`
	HTTPRequest HTTPRequest `mapstructure:"httpRequest"`
	GRPC        GRPC        `mapstructure:"grpc"`
}

func (t TestTrigger) Validate() error {
	switch t.Type {
	case "http":
		if err := t.HTTPRequest.Validate(); err != nil {
			return fmt.Errorf("http request must be valid: %w", err)
		}
	case "grpc":
		if err := t.GRPC.Validate(); err != nil {
			return fmt.Errorf("grpc request must be valid: %w", err)
		}
	case "":
		return fmt.Errorf("type cannot be empty")
	default:
		return fmt.Errorf("type \"%s\" is not supported", t.Type)
	}

	return nil
}

type Output struct {
	Name     string `mapstructure:"name"`
	Selector string `mapstructure:"selector"`
	Value    string `mapstructure:"value"`
}

type TestSpec struct {
	Name       string   `mapstructure:"name"`
	Selector   string   `mapstructure:"selector"`
	Assertions []string `mapstructure:"assertions"`
}
