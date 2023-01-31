import {CallBackProps} from 'react-joyride';

export enum GuidedTours {
  TestRun = 'testRun',
}

export const GuidedTourByPathname = {
  '/test/(.*)/run/(.*)': GuidedTours.TestRun,
};

export interface IGuidedTourState {
  callback: (data: CallBackProps) => void;
  run: boolean;
  stepIndex: number;
  steps: any[];
}
