package definition

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/openapi"
)

type Test struct {
	Id             string           `yaml:"id" json:"id"`
	Name           string           `yaml:"name" json:"name"`
	Description    string           `yaml:"description" json:"description"`
	Trigger        TestTrigger      `yaml:"trigger" json:"trigger"`
	TestDefinition []TestDefinition `yaml:"testDefinition,omitempty" json:"testDefinition,omitempty"`
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
	Type        string              `yaml:"type" json:"type"`
	HTTPRequest openapi.HttpRequest `yaml:"httpRequest" json:"httpRequest"`
	GRPC        openapi.GrpcRequest `yaml:"grpc" json:"grpc"`
}

func (t TestTrigger) Validate() error {
	validTypes := map[string]bool{
		"http": true,
	}

	if t.Type == "" {
		return fmt.Errorf("type cannot be empty")
	}

	if _, ok := validTypes[t.Type]; !ok {
		return fmt.Errorf("type \"%s\" is not supported", t.Type)
	}

	if err := validateRequest(t.HTTPRequest); err != nil {
		return fmt.Errorf("http request must be valid: %w", err)
	}

	return nil
}

type TestDefinition struct {
	Selector   string   `yaml:"selector" json:"selector"`
	Assertions []string `yaml:"assertions" json:"assertions"`
}
