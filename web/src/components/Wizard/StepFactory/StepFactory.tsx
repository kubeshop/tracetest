import {IWizardStep, IWizardStepComponentMap} from 'types/Wizard.types';

const StepComponentMap: IWizardStepComponentMap = {
  Step1: () => <div>Step 1</div>,
  Step2: () => <div>Step 1</div>,
  Step3: () => <div>Step 3</div>,
};

interface IProps {
  step: IWizardStep;
}

const StepFactory = ({step: {component}}: IProps) => {
  const Step = StepComponentMap[component];

  return <Step />;
};

export default StepFactory;
