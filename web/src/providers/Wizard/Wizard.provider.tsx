import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {IWizardState, IWizardStep, TWizardMap} from 'types/Wizard.types';
import Wizard, {isStepEnabled} from 'models/Wizard.model';
import Tracetest from 'redux/apis/Tracetest';

interface IContext extends IWizardState {
  activeStepId: string;
  isLoading: boolean;
  onNext(): void;
  onGoTo(key: string): void;
}

export const Context = createContext<IContext>({
  activeStep: 0,
  activeStepId: '',
  steps: [],
  isLoading: false,
  onNext: noop,
  onGoTo: noop,
});

interface IProps {
  children: React.ReactNode;
  stepsMap: TWizardMap;
}

export const useWizard = () => useContext(Context);

const WizardProvider = ({children, stepsMap}: IProps) => {
  const {useGetWizardQuery, useUpdateWizardMutation} = Tracetest.instance;
  const [updateWizard, {isLoading}] = useUpdateWizardMutation();
  const {data: wizard = Wizard()} = useGetWizardQuery({});
  const steps = useMemo<IWizardStep[]>(
    () =>
      wizard.steps.map((step, index) => ({
        ...step,
        ...(stepsMap[step.id] || {}),
        isEnabled: isStepEnabled(step, index, wizard.steps[index - 1]),
      })),
    [stepsMap, wizard.steps]
  );

  const [activeStep, setActiveStep] = useState(0);

  const activeStepId = steps[activeStep]?.id;
  const isFinalStep = activeStep === steps.length - 1;

  const onNext = useCallback(async () => {
    if (!isFinalStep) {
      await updateWizard({
        steps: wizard.steps.map(step => ({
          ...step,
          ...(step.id === activeStepId ? {state: 'completed'} : {}),
        })),
      });

      setActiveStep(step => step + 1);
    }
  }, [activeStepId, isFinalStep, updateWizard, wizard.steps]);

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
      isLoading,
      onNext,
      onGoTo,
    }),
    [activeStep, activeStepId, isLoading, onGoTo, onNext, steps]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default WizardProvider;
