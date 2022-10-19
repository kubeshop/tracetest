package yaml

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
)

type TestSpecs []TestSpec

func (ts TestSpecs) Model() model.OrderedMap[model.SpanQuery, model.NamedAssertions] {
	mts := model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}
	for _, spec := range ts {
		assertions := make([]model.Assertion, 0, len(spec.Assertions))
		for _, a := range spec.Assertions {
			assertions = append(assertions, model.Assertion(a))
		}

		mts, _ = mts.Add(model.SpanQuery(spec.Selector), model.NamedAssertions{
			Name:       spec.Name,
			Assertions: assertions,
		})
	}
	return mts
}

type Outputs []Output

func (outs Outputs) Model() model.OrderedMap[string, model.Output] {
	mos := model.OrderedMap[string, model.Output]{}
	for _, output := range outs {
		mos, _ = mos.Add(output.Name, model.Output{
			Selector: model.SpanQuery(output.Selector),
			Value:    output.Value,
		})
	}
	return mos
}

type Test struct {
	ID          string      `mapstructure:"id"`
	Name        string      `mapstructure:"name"`
	Description string      `mapstructure:"description"`
	Trigger     TestTrigger `mapstructure:"trigger"`
	Specs       TestSpecs   `mapstructure:"specs,omitempty"`
	Outputs     Outputs     `mapstructure:"outputs,omitempty"`
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

func (t TestTrigger) Model() model.Trigger {
	mt := model.Trigger{
		Type: model.TriggerType(t.Type),
	}

	switch t.Type {
	case "http":
		hr := t.HTTPRequest
		mt.HTTP = &model.HTTPRequest{
			Method:  model.HTTPMethod(hr.Method),
			URL:     hr.URL,
			Headers: hr.Headers.Model(),
			Body:    hr.Body,
			Auth:    hr.Authentication.Model(),
		}
	}
	return mt
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

func (t Test) Model() (model.Test, error) {
	mt := model.Test{
		ID:               id.ID(t.ID),
		Name:             t.Name,
		Description:      t.Description,
		ServiceUnderTest: t.Trigger.Model(),
		Specs:            t.Specs.Model(),
		Outputs:          t.Outputs.Model(),
	}

	return mt, nil
}
