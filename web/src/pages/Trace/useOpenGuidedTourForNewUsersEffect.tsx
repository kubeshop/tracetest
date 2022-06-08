import {Dispatch, SetStateAction, useEffect} from 'react';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';

export function useOpenGuidedTourForNewUsersEffect(setVisible: Dispatch<SetStateAction<boolean>>): void {
  useEffect(() => {
    const isComplete = GuidedTourService.getIsComplete(GuidedTours.Trace);
    if (!isComplete) {
      setVisible(true);
      GuidedTourService.save(GuidedTours.Trace, true);
    }
  }, [setVisible]);
}
