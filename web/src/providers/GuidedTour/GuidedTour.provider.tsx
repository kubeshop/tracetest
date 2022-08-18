import {TourProps, useTour} from '@reactour/tour';
import {createContext, useCallback, useContext, useEffect, useMemo, useState} from 'react';
import {useLocation} from 'react-router-dom';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import HomeStepList from 'components/GuidedTour/homeStepList';
import TraceStepList from 'components/GuidedTour/traceStepList';
import TestDetailsStepList from 'components/GuidedTour/testDetailsStepList';
import {noop} from 'lodash';

interface IContext {
  isTriggerVisible: boolean;
  tour: TourProps;
  setIsTriggerVisible(isVisible: boolean): void;
  onCloseTrigger(): void;
}

export const Context = createContext<IContext>({
  isTriggerVisible: false,
  tour: {} as TourProps,
  setIsTriggerVisible: noop,
  onCloseTrigger: noop,
});

const StepListMap = {
  [GuidedTours.Home]: HomeStepList,
  [GuidedTours.Trace]: TraceStepList,
  [GuidedTours.TestDetails]: TestDetailsStepList,
};

export const useGuidedTour = () => useContext(Context);

const GuidedTourProvider: React.FC = ({children}) => {
  const [isTriggerVisible, setIsTriggerVisible] = useState(false);
  const {pathname} = useLocation();
  const tour = GuidedTourService.getByPathName(pathname);
  const tourProps = useTour();
  const {setSteps} = tourProps;

  useEffect(() => {
    setSteps(StepListMap[tour]);
  }, [tour, setSteps]);

  useEffect(() => {
    const isComplete = GuidedTourService.getIsComplete(tour);
    if (!isComplete) {
      setIsTriggerVisible(true);
    }
  }, [tour]);

  const onCloseTrigger = useCallback(() => {
    setIsTriggerVisible(false);
    GuidedTourService.save(tour, true);
  }, [tour]);

  const value = useMemo(
    () => ({tour: tourProps, isTriggerVisible, setIsTriggerVisible, onCloseTrigger}),
    [tourProps, isTriggerVisible, onCloseTrigger]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default GuidedTourProvider;
