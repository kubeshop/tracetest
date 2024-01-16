package mappings

import (
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/wizard"
)

func (m OpenAPI) Wizard(in *wizard.Wizard) openapi.Wizard {
	steps := make([]openapi.WizardStep, len(in.Steps))
	for i, step := range in.Steps {
		steps[i] = m.WizardStep(step)
	}

	return openapi.Wizard{Steps: steps}
}

func (m OpenAPI) WizardStep(in wizard.Step) openapi.WizardStep {
	step := openapi.WizardStep{
		Id:    in.ID,
		State: string(in.State),
	}

	if in.CompletedAt != nil {
		step.CompletedAt = *in.CompletedAt
	}

	return step
}

func (m Model) Wizard(in openapi.Wizard) wizard.Wizard {
	steps := make([]wizard.Step, len(in.Steps))
	for i, step := range in.Steps {
		steps[i] = m.WizardStep(step)
	}

	return wizard.Wizard{Steps: steps}
}

func (m Model) WizardStep(in openapi.WizardStep) wizard.Step {
	return wizard.Step{
		ID:          in.Id,
		State:       wizard.StepState(in.State),
		CompletedAt: &in.CompletedAt,
	}
}
