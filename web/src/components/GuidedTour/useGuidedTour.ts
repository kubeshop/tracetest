import {useTour} from '@reactour/tour';
import {delay as delayFn} from 'lodash';
import {useEffect, useState} from 'react';
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
  const [isLoaded, setIsLoaded] = useState(false);
  const tourFn = useTour();
  const {setCurrentStep, setIsOpen, setSteps, isOpen} = tourFn;

  useEffect(() => {
    setIsLoaded(false);
    setSteps(StepListMap[tour]);
  }, [tour, setSteps]);

  useEffect(() => {
    const isComplete = GuidedTourService.getIsComplete(tour);
    if (!isComplete) {
      delayFn(() => {
        setCurrentStep(0);
        setIsOpen(true);
        setIsLoaded(true);
      }, delay);
    }
  }, [delay, setCurrentStep, setIsOpen, tour]);

  useEffect(() => {
    if (!isOpen && isLoaded) GuidedTourService.save(tour);
  }, [isLoaded, isOpen, tour]);

  return tourFn;
};

export default useGuidedTour;
