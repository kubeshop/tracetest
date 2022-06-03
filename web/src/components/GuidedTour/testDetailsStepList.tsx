import {StepType} from '@reactour/tour';
import {Typography} from 'antd';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {StepContent} from './StepContent';

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
    content: ({setIsOpen}) => (
      <StepContent title="Test View" setIsOpen={setIsOpen}>
        <Typography.Text>execution time</Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Status),
    content: ({setIsOpen}) => (
      <StepContent title="Test View" setIsOpen={setIsOpen}>
        <Typography.Text>status</Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Assertions),
    content: ({setIsOpen}) => (
      <StepContent title="Test View" setIsOpen={setIsOpen}>
        <Typography.Text>assertions</Typography.Text>
      </StepContent>
    ),
    highlightedSelectors: [
      GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Assertions),
      GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Passed),
      GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.Failed),
    ],
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.TestDetails, Steps.RunTest),
    content: ({setIsOpen}) => (
      <StepContent title="Test View" setIsOpen={setIsOpen}>
        <Typography.Text>assertions</Typography.Text>
      </StepContent>
    ),
  },
];

export default StepList;
