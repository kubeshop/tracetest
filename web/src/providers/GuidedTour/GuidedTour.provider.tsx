import {TourProps, useTour} from '@reactour/tour';
import {createContext, useContext, useEffect, useMemo, useState} from 'react';
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
}

export const Context = createContext<IContext>({
  isTriggerVisible: false,
  tour: {} as TourProps,
  setIsTriggerVisible: noop,
});

const StepListMap = {
  [GuidedTours.Home]: HomeStepList,
  [GuidedTours.Trace]: TraceStepList,
  [GuidedTours.TestDetails]: TestDetailsStepList,
};

export const useGuidedTour = () => useContext(Context);

const GuidedTourProvider: React.FC = ({children}) => {
  const [isTriggerVisible, setIsTriggerVisible] = useState(false);
  const [isLoaded, setIsLoaded] = useState(false);
  const {pathname} = useLocation();

  const tour = GuidedTourService.getByPathName(pathname);

  const tourProps = useTour();
  const {setCurrentStep, setSteps, isOpen} = tourProps;

  useEffect(() => {
    setIsLoaded(false);
    setSteps(StepListMap[tour]);
  }, [tour, setSteps]);

  useEffect(() => {
    if (!isOpen && isLoaded) {
      GuidedTourService.save(tour);
      setCurrentStep(0);
    }
  }, [isLoaded, isOpen, setCurrentStep, tour]);

  useEffect(() => {
    const isComplete = GuidedTourService.getIsComplete(tour);
    if (!isComplete) {
      setIsTriggerVisible(true);
      GuidedTourService.save(tour, true);
    }
  }, [tour]);

  const value = useMemo(
    () => ({tour: tourProps, isTriggerVisible, setIsTriggerVisible}),
    [tourProps, isTriggerVisible]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default GuidedTourProvider;
