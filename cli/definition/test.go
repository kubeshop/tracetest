package definition

import "fmt"

type Test struct {
	Id          string      `yaml:"id"`
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Trigger     TestTrigger `yaml:"trigger"`
	Spec        []TestSpec  `yaml:"specs,omitempty"`
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
	Type        string      `yaml:"type"`
	HTTPRequest HttpRequest `yaml:"httpRequest"`
	GRPC        GrpcRequest `yaml:"grpc"`
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

type TestSpec struct {
	Selector   string   `yaml:"selector"`
	Assertions []string `yaml:"assertions"`
}
