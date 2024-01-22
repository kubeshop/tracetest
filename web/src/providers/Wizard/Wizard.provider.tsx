import {noop} from 'lodash';
import Wizard, {TWizardStepId, isStepEnabled} from 'models/Wizard.model';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import Tracetest from 'redux/apis/Tracetest';
import WizardAnalytics from 'services/Analytics/Wizard.service';
import {IWizardState, IWizardStep, TWizardMap} from 'types/Wizard.types';
import WizardService from 'services/Wizard.service';

interface IContext extends IWizardState {
  activeStepId: string;
  isLoading: boolean;
  onNext(): void;
  onGoTo(key: string): void;
  onCompleteStep(stepId: TWizardStepId): void;
}

export const Context = createContext<IContext>({
  activeStep: 0,
  activeStepId: '',
  steps: [],
  isLoading: false,
  onNext: noop,
  onGoTo: noop,
  onCompleteStep: noop,
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
      wizard.steps
        .filter(step => !!stepsMap[step.id])
        .map((step, index) => ({
          ...step,
          ...(stepsMap[step.id]! || {}),
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
    WizardAnalytics.onStepComplete(activeStepId);
  }, [activeStepId, isFinalStep, updateWizard, wizard.steps]);

  const onCompleteStep = useCallback(
    async (stepId: TWizardStepId) => {
      if (WizardService.shouldUpdate(stepId, wizard.steps)) {
        const updatedSteps = WizardService.completeStep(stepId, wizard.steps);
        await updateWizard({steps: updatedSteps}).unwrap();
        WizardAnalytics.onStepComplete(stepId);
      }
    },
    [updateWizard, wizard.steps]
  );

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
      onCompleteStep,
    }),
    [activeStep, activeStepId, isLoading, onCompleteStep, onGoTo, onNext, steps]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default WizardProvider;
