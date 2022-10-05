import {Dispatch, SetStateAction, useEffect} from 'react';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {OnboardingState} from './GuidedTour.provider';

export function useShowOnboardingWhenNotCompletedEffect(
  tour: GuidedTours,
  setState: Dispatch<SetStateAction<OnboardingState>>
) {
  useEffect(() => {
    if (tour === GuidedTours.Trace) {
      const isComplete = GuidedTourService.getIsComplete(tour);
      if (!isComplete) {
        setState(st => ({...st, dialog: true}));
      }
    }
  }, [setState, tour]);
  return tour;
}
