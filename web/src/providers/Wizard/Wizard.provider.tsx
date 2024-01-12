import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {IWizardState, IWizardStep} from 'types/Wizard.types';

interface IContext extends IWizardState {
  activeStepId: string;
  isLoading: boolean;
  onNext(): void;
  onPrev(): void;
  onGoTo(key: string): void;
}

export const Context = createContext<IContext>({
  activeStep: 0,
  activeStepId: '',
  steps: [],
  isLoading: false,
  onNext: noop,
  onPrev: noop,
  onGoTo: noop,
});

interface IProps {
  children: React.ReactNode;
  steps: IWizardStep[];
}

export const useWizard = () => useContext(Context);

const WizardProvider = ({children, steps = []}: IProps) => {
  const [activeStep, setActiveStep] = useState(0);

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

  const onGoTo = useCallback(
    key => {
      const index = steps.findIndex(step => step.id === key);
      setActiveStep(index);
    },
    [steps]
  );

  const value = useMemo<IContext>(
    () => ({
      activeStep,
      activeStepId,
      steps,
      isLoading: false, // TODO: implement loading
      onNext,
      onPrev,
      onGoTo,
    }),
    [activeStep, activeStepId, onGoTo, onNext, onPrev, steps]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default WizardProvider;
