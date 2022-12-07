import {noop} from 'lodash';
import React, {createContext, Dispatch, SetStateAction, useCallback, useContext, useMemo, useState} from 'react';
import Joyride, {CallBackProps} from 'react-joyride';
import {useAppDispatch} from '../../redux/hooks';
import {setUserPreference} from '../../redux/slices/User.slice';
import GuidedTourService from '../../services/GuidedTour.service';
import {useGetTooltipComponent} from './useGetTooltipComponent';
import {useShowOnboardingWhenNotCompletedEffect} from './useShowOnboardingWhenNotCompletedEffect';
import {useUpdateStepsBasedOnLocationEffect} from './useUpdateStepsBasedOnLocationEffect';

interface IContext {
  state: OnboardingState;
  setState: Dispatch<SetStateAction<OnboardingState>>;
  onSkip(): void;
}

export const Context = createContext<IContext>({
  state: {} as OnboardingState,
  setState: noop,
  onSkip: noop,
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
  const dispatch = useAppDispatch();
  const tour = GuidedTourService.useGetCurrentOnboardingLocation();
  useUpdateStepsBasedOnLocationEffect(state, setState);
  useShowOnboardingWhenNotCompletedEffect(tour, setState);

  const onSkip = useCallback(() => {
    dispatch(
      setUserPreference({
        key: 'isOnboardingComplete',
        value: true,
      })
    );

    setState(st => ({...st, dialog: false}));
  }, [dispatch]);

  const value = useMemo(() => ({state, setState, onSkip}), [setState, state, onSkip]);

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
