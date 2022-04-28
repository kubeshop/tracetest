import {StepType} from '@reactour/tour';
import GuidedTourService, {GuidedTours} from '../../entities/GuidedTour/GuidedTour.service';

export enum Steps {
  Diagram = 'diagram',
  SpanDetail = 'spanDetail',
  TestResults = 'testResults',
  Timeline = 'timeline',
}

const StepList: StepType[] = [
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.Diagram),
    content: 'diagram',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.SpanDetail),
    content: 'span detail',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.TestResults),
    content: 'test results',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.Timeline),
    content: 'timeline',
  },
];

export default StepList;
