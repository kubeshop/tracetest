import {Dispatch, SetStateAction} from 'react';
import {ACTIONS, CallBackProps, EVENTS, STATUS, Step} from 'react-joyride';

import {GuidedTours, IGuidedTourState} from 'constants/GuidedTour';
import TestRunSteps, {switchTestRunMode} from './testRunSteps';

export const StepsByGuidedTour: Record<GuidedTours, Step[]> = {
  [GuidedTours.TestRun]: TestRunSteps,
};

export const CallbackByGuidedTour: Record<
  GuidedTours,
  (setState: Dispatch<SetStateAction<IGuidedTourState>>) => (data: CallBackProps) => void
> = {
  [GuidedTours.TestRun]: (setState: Dispatch<SetStateAction<IGuidedTourState>>) => (data: CallBackProps) => {
    const {action, index, status, type} = data;

    if (([STATUS.FINISHED, STATUS.SKIPPED] as string[]).includes(status)) {
      // Need to set our running state to false, so we can restart if we click start again.
      setState(st => ({...st, run: false, stepIndex: 0}));
      return;
    }

    if (([EVENTS.STEP_AFTER] as string[]).includes(type)) {
      const nextStepIndex = index + (action === ACTIONS.PREV ? -1 : 1);

      if (action === ACTIONS.NEXT && index === 2) {
        // Need to switch the page to Test mode. The Test mode page will update the step on mount.
        switchTestRunMode(2);
      }

      // Move to the next step.
      setState(st => ({...st, stepIndex: nextStepIndex}));

      if (action === ACTIONS.PREV && index === 2) {
        // Need to switch the page to Trigger mode.
        switchTestRunMode(0);
      }
    }
  },
};
