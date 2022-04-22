import {StepType} from '@reactour/tour';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTourService';

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
    content: 'Create test',
    disableActions: true,
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Method),
    content: 'method',
    disableActions: false,
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Url),
    content: 'url',
    disableActions: false,
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Name),
    content: 'name',
    disableActions: false,
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Headers),
    content: 'headers',
    disableActions: false,
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Body),
    content: 'body',
    disableActions: false,
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Run),
    content: 'run',
    disableActions: false,
  },
];

export default StepList;
