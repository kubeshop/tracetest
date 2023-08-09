package yaml

import (
	"fmt"

	dc "github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/testsuite"
)

type TestSuite struct {
	ID          string            `mapstructure:"id"`
	Name        string            `mapstructure:"name"`
	Description string            `mapstructure:"description" yaml:",omitempty"`
	Env         map[string]string `mapstructure:"env" yaml:",omitempty"`
	Steps       []string          `mapstructure:"steps"`
}

func (t TestSuite) Model() testsuite.TestSuite {
	mt := testsuite.TestSuite{}
	dc.DeepCopy(t, &mt)
	mt.StepIDs = make([]id.ID, 0, len(t.Steps))
	for _, stepID := range t.Steps {
		mt.StepIDs = append(mt.StepIDs, id.ID(stepID))
	}

	return mt
}

func (t TestSuite) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("suite name cannot be empty")
	}

	return nil
}
