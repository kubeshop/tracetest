import {TWizardStepId, WizardStep} from 'models/Wizard.model';

export interface IWizardState {
  activeStep: number;
  steps: IWizardStep[];
}

export interface IWizardStepMetadata {
  name: string;
  description: string;
  component: React.FC;
}

export interface IWizardStep extends WizardStep, IWizardStepMetadata {}
export type TWizardMap = Record<TWizardStepId, IWizardStepMetadata>;

interface IWizardStepComponentProps {}

export interface IWizardStepComponentMap
  extends Record<string, (props: IWizardStepComponentProps) => React.ReactElement> {}
