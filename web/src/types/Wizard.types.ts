import {TWizardStepId, WizardStep} from 'models/Wizard.model';

export interface IWizardState {
  activeStep: number;
  steps: IWizardStep[];
}

export interface IWizardStepMetadata {
  name: string;
  description: string;
  component: React.FC<IWizardStepComponentProps>;
  tabComponent: React.FC;
  isEnabled: boolean;
}
export interface IWizardStepComponentProps {
  isLoading: boolean;
  onNext(): void;
}

export interface IWizardStep extends WizardStep, IWizardStepMetadata {}
export type TWizardMap = Record<TWizardStepId, IWizardStepMetadata>;
