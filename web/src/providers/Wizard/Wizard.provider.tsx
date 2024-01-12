import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {IWizardState, IWizardStep} from 'types/Wizard.types';

interface IContext extends IWizardState {
  activeStepId: string;
  isLoading: boolean;
  onNext(): void;
  onPrev(): void;
}

export const Context = createContext<IContext>({
  activeStep: 0,
  activeStepId: '',
  steps: [],
  isLoading: false,
  onNext: noop,
  onPrev: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useWizard = () => useContext(Context);

const initialSteps: IWizardStep[] = [
  {
    id: 'step1',
    name: 'Step 1',
    description: 'Step 1 description',
    component: 'Step1',
  },
  {
    id: 'step2',
    name: 'Step 2',
    description: 'Step 2 description',
    component: 'Step2',
  },
  {
    id: 'step3',
    name: 'Step 3',
    description: 'Step 3 description',
    component: 'Step3',
  },
];

const WizardProvider = ({children}: IProps) => {
  const [activeStep, setActiveStep] = useState(0);
  const [steps, setSteps] = useState<IWizardStep[]>(initialSteps);

  const activeStepId = steps[activeStep]?.id;
  const isFinalStep = activeStep === steps.length - 1;

  const onNext = useCallback(() => {
    if (!isFinalStep) {
      setActiveStep(step => step + 1);
    }
  }, [isFinalStep]);

  const onPrev = useCallback(() => {
    setActiveStep(step => step - 1);
  }, []);

  const value = useMemo<IContext>(
    () => ({
      activeStep,
      activeStepId,
      steps,
      isLoading: false, // TODO: implement loading
      onNext,
      onPrev,
    }),
    [activeStep, activeStepId, onNext, onPrev, steps]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default WizardProvider;
