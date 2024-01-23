import {TWizardStepId, isStepCompleted, WizardStep} from 'models/Wizard.model';

const WizardService = () => ({
  shouldUpdate(id: TWizardStepId, steps: WizardStep[]): boolean {
    const step = steps.find(({id: stepId}) => stepId === id);

    return !!step && !isStepCompleted(step);
  },
  completeStep(id: TWizardStepId, steps: WizardStep[]): WizardStep[] {
    const stepIndex = steps.findIndex(({id: stepId}) => stepId === id);

    if (stepIndex === -1) {
      return steps;
    }

    return [
      ...steps.slice(0, stepIndex + 1).map(step => {
        if (isStepCompleted(step)) {
          return step;
        }

        const update: WizardStep = {
          id: step.id,
          state: 'completed',
          completedAt: new Date().toISOString(),
        };

        return update;
      }),
      ...steps.slice(stepIndex + 1),
    ];
  },
});

export default WizardService();
