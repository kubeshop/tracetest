import {StepType} from '@reactour/tour';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';

export enum Steps {
  ExecutionTime = 'executionTime',
  Status = 'status',
  Assertions = 'assertions',
  RunTest = 'runTest',
  Passed = 'passed',
  Failed = 'failed',
}

const StepList: StepType[] = [
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.ExecutionTime),
    content: 'execution time',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Status),
    content: 'status',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Assertions),
    content: 'assertions',
    highlightedSelectors: [
      GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Assertions),
      GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Passed),
      GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Failed),
    ],
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.RunTest),
    content: 'run test',
  },
];

export default StepList;
