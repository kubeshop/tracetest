import {Step} from 'react-joyride';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';

export enum Steps {
  ExecutionTime = 'executionTime',
  Status = 'status',
  Assertions = 'assertions',
  RunTest = 'runTest',
  Passed = 'passed',
  Failed = 'failed',
}

const StepList: Step[] = [
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.ExecutionTime),
    content: 'execution time',
  },
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Status),
    content: 'status',
  },
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Assertions),
    content: 'assertions',
  },
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.RunTest),
    content: 'assertions',
  },
];

export default StepList;
