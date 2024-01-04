import {noop} from 'lodash';
import React, {createContext, useCallback, useContext, useMemo, useState} from 'react';
import Joyride from 'react-joyride';
import {useLocation} from 'react-router-dom';

import StepContent from 'components/GuidedTour/StepContent';
import {CallbackByGuidedTour, StepsByGuidedTour} from 'components/GuidedTour/steps';
import {switchTestRunMode} from 'components/GuidedTour/testRunSteps';
import {IGuidedTourState} from 'constants/GuidedTour';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import GuidedTourService from 'services/GuidedTour.service';

const {onGuidedTourClick} = HomeAnalyticsService;

interface IContext {
  isGuidedTourRunning: boolean;
  onStart(): void;
  setGuidedTourStep(index: number): void;
}

const Context = createContext<IContext>({
  isGuidedTourRunning: false,
  onStart: noop,
  setGuidedTourStep: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useGuidedTour = () => useContext(Context);

const GuidedTourProvider = ({children}: IProps) => {
  const pathname = useLocation().pathname;
  const tourByPathname = GuidedTourService.getByPathName(pathname);
  const [guidedTourState, setGuidedTourState] = useState<IGuidedTourState>({
    callback: noop,
    run: false,
    stepIndex: 0,
    steps: [],
  });

  const onStart = useCallback(() => {
    if (!tourByPathname) return;

    onGuidedTourClick(); // Analytics
    switchTestRunMode(0);
    setGuidedTourState(prev => ({
      ...prev,
      callback: CallbackByGuidedTour[tourByPathname](setGuidedTourState),
      run: true,
      stepIndex: 0,
      steps: StepsByGuidedTour[tourByPathname],
    }));
  }, [tourByPathname]);

  const setGuidedTourStep = useCallback((stepIndex: number) => {
    setGuidedTourState(prev => ({...prev, stepIndex}));
  }, []);

  const value = useMemo(
    () => ({isGuidedTourRunning: guidedTourState.run, onStart, setGuidedTourStep}),
    [guidedTourState.run, onStart, setGuidedTourStep]
  );

  return (
    <Context.Provider value={value}>
      <Joyride
        callback={guidedTourState.callback}
        continuous
        disableOverlay
        run={guidedTourState.run}
        stepIndex={guidedTourState.stepIndex}
        steps={guidedTourState.steps}
        tooltipComponent={StepContent}
      />
      {children}
    </Context.Provider>
  );
};

export default GuidedTourProvider;
