import {noop} from 'lodash';
import React, {createContext, Dispatch, SetStateAction, useContext, useMemo, useState} from 'react';
import Joyride, {CallBackProps} from 'react-joyride';
import GuidedTourService from '../../services/GuidedTour.service';
import {useGetTooltipComponent} from './useGetTooltipComponent';
import {useShowOnboardingWhenNotCompletedEffect} from './useShowOnboardingWhenNotCompletedEffect';
import {useUpdateStepsBasedOnLocationEffect} from './useUpdateStepsBasedOnLocationEffect';

interface IContext {
  state: OnboardingState;
  setState: Dispatch<SetStateAction<OnboardingState>>;
}

export const Context = createContext<IContext>({
  state: {} as OnboardingState,
  setState: noop,
});

export interface OnboardingState {
  callback: (data: CallBackProps) => void;
  stepIndex: number;
  tourActive: boolean;
  run: boolean;
  dialog: boolean;
  steps: any[];
}

export const useGuidedTour = () => useContext(Context);

const GuidedTourProvider: React.FC = ({children}) => {
  const [state, setState] = useState<OnboardingState>({
    run: false,
    dialog: false,
    stepIndex: 0,
    steps: [],
    tourActive: false,
    callback: noop,
  });
  const tour = GuidedTourService.useGetCurrentOnboardingLocation();
  useUpdateStepsBasedOnLocationEffect(state, setState);
  useShowOnboardingWhenNotCompletedEffect(tour, setState);
  const value = useMemo(() => ({state, setState}), [setState, state]);
  return (
    <Context.Provider value={value}>
      <Joyride
        disableOverlay
        callback={state.callback}
        continuous
        run={state.run}
        stepIndex={state.stepIndex}
        steps={state.steps}
        tooltipComponent={useGetTooltipComponent(tour)}
      />
      {children}
    </Context.Provider>
  );
};

export default GuidedTourProvider;
