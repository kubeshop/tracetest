import {Button, notification, Space} from 'antd';
import {noop} from 'lodash';
import React, {createContext, useCallback, useContext, useEffect, useMemo, useState} from 'react';
import Joyride from 'react-joyride';
import {useLocation} from 'react-router-dom';

import StepContent from 'components/GuidedTour/StepContent';
import {CallbackByGuidedTour, StepsByGuidedTour} from 'components/GuidedTour/steps';
import {switchTestRunMode} from 'components/GuidedTour/testRunSteps';
import {IGuidedTourState} from 'constants/GuidedTour';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setUserPreference} from 'redux/slices/User.slice';
import UserSelectors from 'selectors/User.selectors';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import GuidedTourService from 'services/GuidedTour.service';
import {useNotification} from '../Notification/Notification.provider';

const {onGuidedTourClick} = HomeAnalyticsService;
const NOTIFICATION_KEY = 'guided-tour-notification';

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
  const dispatch = useAppDispatch();
  const pathname = useLocation().pathname;
  const tourByPathname = GuidedTourService.getByPathName(pathname);
  const {showNotification} = useNotification();
  const [guidedTourState, setGuidedTourState] = useState<IGuidedTourState>({
    callback: noop,
    run: false,
    stepIndex: 0,
    steps: [],
  });
  const showGuidedTourNotification = useAppSelector(state =>
    UserSelectors.selectUserPreference(state, 'showGuidedTourNotification')
  );

  useEffect(() => {
    if (!tourByPathname || !showGuidedTourNotification) return;

    const onSkip = () => {
      notification.close(NOTIFICATION_KEY);
      dispatch(setUserPreference({key: 'showGuidedTourNotification', value: false}));
    };

    const onStart = () => {
      onGuidedTourClick(); // Analytics
      onSkip();
      switchTestRunMode(0);
      setGuidedTourState(prev => ({
        ...prev,
        callback: CallbackByGuidedTour[tourByPathname](setGuidedTourState),
        run: true,
        stepIndex: 0,
        steps: StepsByGuidedTour[tourByPathname],
      }));
    };

    const btn = (
      <Space>
        <Button data-cy="guided-tour-cancel-notification" ghost onClick={onSkip} type="primary">
          No thanks
        </Button>
        <Button onClick={onStart} type="primary">
          Show me around
        </Button>
      </Space>
    );

    showNotification({
      type: 'open',
      description: 'Walk through Tracetest features',
      duration: 0,
      btn,
      key: NOTIFICATION_KEY,
      title: 'Do you want to take a quick tour of Tracetest?',
      onClose: onSkip,
    });
  }, [dispatch, showGuidedTourNotification, showNotification, tourByPathname]);

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
