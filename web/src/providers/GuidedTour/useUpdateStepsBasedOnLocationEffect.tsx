import {Dispatch, SetStateAction, useEffect} from 'react';
import {CallBackProps, Step} from 'react-joyride';
import {useLocation} from 'react-router-dom';
import HomeStepList from '../../components/GuidedTour/homeStepList';
import TestDetailsStepList from '../../components/GuidedTour/testDetailsStepList';
import TraceStepList, {switchTraceMode} from '../../components/GuidedTour/traceStepList';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {OnboardingState} from './GuidedTour.provider';

const StepListMap: Record<GuidedTours, Step[]> = {
  [GuidedTours.Home]: HomeStepList,
  [GuidedTours.Trace]: TraceStepList,
  [GuidedTours.TestDetails]: TestDetailsStepList,
};
const CallbackListMap: Record<
  GuidedTours,
  (setState: Dispatch<SetStateAction<OnboardingState>>) => (data: CallBackProps) => void
> = {
  [GuidedTours.Home]: () => () => {},
  [GuidedTours.Trace]: (setState: Dispatch<SetStateAction<OnboardingState>>) => (data: CallBackProps) => {
    const {action, index, type} = data;
    if (type === 'tour:end' || action === 'close') {
      setState(st => ({...st, run: false, stepIndex: 0, tourActive: false}));
      return;
    }
    if (type === 'step:after') {
      if (action === 'prev') {
        setState(st => ({...st, stepIndex: st.stepIndex - 1}));
        if (type === 'step:after' && index === 2) {
          switchTraceMode(0)();
        }
        return;
      }
      if (index === 2 /* or step.target === '#home' */) {
        setState(st => ({...st, run: false}));
        switchTraceMode(2)();
      }
      setState(st => ({...st, stepIndex: st.stepIndex + 1}));
    }
  },
  [GuidedTours.TestDetails]: () => () => {},
};

export function useUpdateStepsBasedOnLocationEffect(
  state: OnboardingState,
  setState: Dispatch<SetStateAction<OnboardingState>>
) {
  const {pathname} = useLocation();
  useEffect(() => {
    let byPathName = GuidedTourService.getByPathName(pathname);
    return setState(st => ({
      ...st,
      steps: StepListMap[byPathName],
      callback: CallbackListMap[byPathName](setState),
    }));
  }, [pathname, setState]);
}
