package yaml

import (
	"fmt"

	dc "github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/test"
)

type TestSpecs []TestSpec

func (ts TestSpecs) Model() test.Specs {
	specs := make(test.Specs, 0, len(ts))
	for _, spec := range ts {
		assertions := make([]test.Assertion, 0, len(spec.Assertions))
		for _, a := range spec.Assertions {
			assertions = append(assertions, test.Assertion(a))
		}

		specs = append(specs, test.TestSpec{
			Name:       spec.Name,
			Selector:   test.SpanQuery(spec.Selector),
			Assertions: assertions,
		})
	}
	return specs
}

type Outputs []Output

func (outs Outputs) Model() test.Outputs {
	outputs := make(test.Outputs, 0, len(outs))
	for _, output := range outs {
		outputs = append(outputs, test.Output{
			Name:     output.Name,
			Selector: test.SpanQuery(output.Selector),
			Value:    output.Value,
		})
	}
	return outputs
}

type Test struct {
	ID          string      `mapstructure:"id"`
	Name        string      `mapstructure:"name"`
	Description string      `mapstructure:"description" yaml:",omitempty"`
	Trigger     TestTrigger `mapstructure:"trigger" dc:"serviceUnderTest"`
	Specs       TestSpecs   `mapstructure:"specs" yaml:",omitempty"`
	Outputs     Outputs     `mapstructure:"outputs,omitempty" yaml:",omitempty"`
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
	Type        string         `mapstructure:"type"`
	HTTPRequest HTTPRequest    `mapstructure:"httpRequest" yaml:"httpRequest,omitempty" dc:"http"`
	GRPC        GRPC           `mapstructure:"grpc" yaml:"grpc,omitempty"`
	TRACEID     TRACEIDRequest `mapstructure:"traceid" yaml:"traceid,omitempty"`
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
	case "traceid":
		if err := t.TRACEID.Validate(); err != nil {
			return fmt.Errorf("traceid request must be valid: %w", err)
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
	Name       string   `mapstructure:"name" yaml:",omitempty"`
	Selector   string   `mapstructure:"selector"`
	Assertions []string `mapstructure:"assertions"`
}

func (t Test) Model() test.Test {
	mt := test.Test{}
	dc.DeepCopy(t, &mt)
	mt.Specs = t.Specs.Model()
	mt.Outputs = t.Outputs.Model()

	return mt
}
