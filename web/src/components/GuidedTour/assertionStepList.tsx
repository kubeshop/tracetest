import {StepType} from '@reactour/tour';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';

export enum Steps {
  Selectors = 'selectors',
  Checks = 'checks',
}

const StepList: StepType[] = [
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Assertion, Steps.Selectors),
    content: 'selectors',
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Assertion, Steps.Checks),
    content: 'checks',
  },
];

export default StepList;
