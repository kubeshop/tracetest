package definition

import (
	"fmt"
)

type Test struct {
	Id          string      `yaml:"id" json:"id"`
	Name        string      `yaml:"name" json:"name"`
	Description string      `yaml:"description" json:"description"`
	Trigger     TestTrigger `yaml:"trigger" json:"trigger"`
	Specs       []TestSpec  `yaml:"specs,omitempty" json:"specs,omitempty"`
	Outputs     []Output    `yaml:"outputs,omitempty" json:"outputs,omitempty"`
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
	Type        string      `yaml:"type" json:"type"`
	HTTPRequest HTTPRequest `yaml:"httpRequest" json:"httpRequest"`
	GRPC        GRPC        `yaml:"grpc" json:"grpc"`
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
	Name     string `yaml:"name"`
	Selector string `yaml:"selector"`
	Value    string `yaml:"value"`
}

type TestSpec struct {
	Name       string   `yaml:"name" json:"name"`
	Selector   string   `yaml:"selector" json:"selector"`
	Assertions []string `yaml:"assertions" json:"assertions"`
}
