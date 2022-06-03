import {StepType} from '@reactour/tour';
import {Typography} from 'antd';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {StepContent} from './StepContent';

export enum Steps {
  CreateTest = 'create_test',
  Method = 'method',
  Url = 'url',
  Name = 'name',
  Headers = 'headers',
  Body = 'body',
  Run = 'run',
}

const StepList: StepType[] = [
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.CreateTest),
    content: ({setIsOpen}) => (
      <StepContent title="Create test" setIsOpen={setIsOpen}>
        <Typography.Text>Create test</Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Method),
    content: ({setIsOpen}) => (
      <StepContent title="method" setIsOpen={setIsOpen}>
        <Typography.Text>method</Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Url),
    content: ({setIsOpen}) => (
      <StepContent title="url" setIsOpen={setIsOpen}>
        <Typography.Text>url</Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Name),
    content: ({setIsOpen}) => (
      <StepContent title="name" setIsOpen={setIsOpen}>
        <Typography.Text>name</Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Headers),
    content: ({setIsOpen}) => (
      <StepContent title="headers" setIsOpen={setIsOpen}>
        <Typography.Text>headers</Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Body),
    content: ({setIsOpen}) => (
      <StepContent title="body" setIsOpen={setIsOpen}>
        <Typography.Text>body</Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Run),
    content: ({setIsOpen}) => (
      <StepContent title="Test View" setIsOpen={setIsOpen}>
        <Typography.Text>run</Typography.Text>
      </StepContent>
    ),
  },
];

export default StepList;
