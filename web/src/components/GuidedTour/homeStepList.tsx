import {StepType} from '@reactour/tour';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';

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
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Method),
    content: 'method',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Url),
    content: 'url',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Name),
    content: 'name',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Headers),
    content: 'headers',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Body),
    content: 'body',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Run),
    content: 'run',
  },
];

export default StepList;
