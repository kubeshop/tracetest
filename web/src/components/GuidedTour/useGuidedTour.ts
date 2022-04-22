import {useTour} from '@reactour/tour';
import {delay as delayFn} from 'lodash';
import {useEffect} from 'react';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTourService';
import HomeStepList from './homeStepList';
import AssertionStepList from './assertionStepList';
import TraceStepList from './traceStepList';
import TestDetailsStepList from './testDetailsStepList';

const StepListMap = {
  [GuidedTours.Home]: HomeStepList,
  [GuidedTours.Assertion]: AssertionStepList,
  [GuidedTours.Trace]: TraceStepList,
  [GuidedTours.TestDetails]: TestDetailsStepList,
};

const useGuidedTour = (tour: GuidedTours, delay = 500) => {
  const tourFn = useTour();
  const {setCurrentStep, setIsOpen, setSteps} = tourFn;

  useEffect(() => {
    setSteps(StepListMap[tour]);
  }, [tour, setSteps]);

  useEffect(() => {
    const isComplete = GuidedTourService.getIsComplete(tour);
    if (!isComplete) {
      delayFn(() => {
        setCurrentStep(0);
        setIsOpen(true);
      }, delay);
    }
  }, [delay, setCurrentStep, setIsOpen, tour]);

  return tourFn;
};

export default useGuidedTour;
