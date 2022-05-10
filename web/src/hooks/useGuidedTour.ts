import {useTour} from '@reactour/tour';
import {useEffect, useState} from 'react';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import HomeStepList from 'components/GuidedTour/homeStepList';
import AssertionStepList from 'components/GuidedTour/assertionStepList';
import TraceStepList from 'components/GuidedTour/traceStepList';
import TestDetailsStepList from 'components/GuidedTour/testDetailsStepList';

const StepListMap = {
  [GuidedTours.Home]: HomeStepList,
  [GuidedTours.Assertion]: AssertionStepList,
  [GuidedTours.Trace]: TraceStepList,
  [GuidedTours.TestDetails]: TestDetailsStepList,
};

const useGuidedTour = (tour: GuidedTours) => {
  const [isLoaded, setIsLoaded] = useState(false);
  const tourFn = useTour();
  const {setCurrentStep, setSteps, isOpen} = tourFn;

  useEffect(() => {
    setIsLoaded(false);
    setSteps(StepListMap[tour]);
  }, [tour, setSteps]);

  // useEffect(() => {
  //   const isComplete = GuidedTourService.getIsComplete(tour);
  //   if (!isComplete) {
  //     delayFn(() => {
  //       setCurrentStep(0);
  //       setIsOpen(true);
  //       setIsLoaded(true);
  //     }, delay);
  //   }
  // }, [delay, setCurrentStep, setIsOpen, tour]);

  useEffect(() => {
    if (!isOpen && isLoaded) {
      GuidedTourService.save(tour);
      setCurrentStep(0);
    }
  }, [isLoaded, isOpen, setCurrentStep, tour]);

  return tourFn;
};

export default useGuidedTour;
