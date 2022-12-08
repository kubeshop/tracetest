import {Dispatch, SetStateAction, useEffect} from 'react';
import {useAppSelector} from '../../redux/hooks';
import UserSelectors from '../../selectors/User.selectors';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {OnboardingState} from './GuidedTour.provider';

export function useShowOnboardingWhenNotCompletedEffect(
  tour: GuidedTours,
  setState: Dispatch<SetStateAction<OnboardingState>>
) {
  const isUserComplete = useAppSelector(state => UserSelectors.selectUserPreference(state, 'isOnboardingComplete'));

  useEffect(() => {
    if (tour === GuidedTours.Trace) {
      const isComplete = GuidedTourService.getIsComplete(tour);
      if (!isComplete && !isUserComplete) {
        setState(st => ({...st, dialog: true}));
      }
    }
  }, [isUserComplete, setState, tour]);
  return tour;
}
