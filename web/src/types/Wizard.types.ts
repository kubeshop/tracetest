export interface IWizardState {
  activeStep: number;
  steps: IWizardStep[];
}

export interface IWizardStep {
  id: string;
  name: string;
  description: string;
  component: string; // enum?
  status?: TWizardStepStatus;
}

type TWizardStepStatus = 'complete' | 'pending';

interface IWizardStepComponentProps {}

export interface IWizardStepComponentMap
  extends Record<string, (props: IWizardStepComponentProps) => React.ReactElement> {}
