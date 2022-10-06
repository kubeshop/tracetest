import {Step} from 'react-joyride';
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

const StepList: Step[] = [
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.CreateTest),
    content: 'create test',
  },
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Method),
    content: 'method',
  },
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Url),
    content: 'url',
  },
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Name),
    content: 'name',
  },
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Headers),
    content: 'headers',
  },
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Body),
    content: 'body',
  },
  {
    target: GuidedTourService.getSelectorStep(GuidedTours.Home, Steps.Run),
    content: 'run',
  },
];

export default StepList;
