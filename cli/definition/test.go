package definition

import "fmt"

type Test struct {
	Id          string      `yaml:"id"`
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Trigger     TestTrigger `yaml:"trigger"`
}

func (t Test) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("test name cannot be empty")
	}

	if t.Description == "" {
		return fmt.Errorf("test description cannot be empty")
	}

	if err := t.Trigger.Validate(); err != nil {
		return fmt.Errorf("test trigger must be valid: %w", err)
	}

	return nil
}

type TestTrigger struct {
	Type        string      `yaml:"type"`
	HTTPRequest HttpRequest `yaml:"http_request"`
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

	if err := t.HTTPRequest.Validate(); err != nil {
		return fmt.Errorf("http request must be valid: %w", err)
	}

	return nil
}

type TestDefinition struct {
	Selector   string   `yaml:"selector"`
	Assertions []string `yaml:"assertions"`
}
