import {Model, TWizardSchemas} from 'types/Common.types';

export type TRawWizard = TWizardSchemas['Wizard'];
export type TRawWizardStep = TWizardSchemas['WizardStep'];
export type TWizardStepId = NonNullable<TRawWizardStep['id']>;

type Wizard = Model<TRawWizard, {steps: WizardStep[]}>;
export type WizardStep = Model<TRawWizardStep, {id: TWizardStepId}>;

const defaultWizard: TRawWizard = {steps: []};

function Wizard({steps = []} = defaultWizard): Wizard {
  return {
    steps: steps.map(({id = 'tracing_backend', state = 'pending', completedAt = ''}) => ({
      id,
      state,
      completedAt,
    })),
  };
}

export default Wizard;
