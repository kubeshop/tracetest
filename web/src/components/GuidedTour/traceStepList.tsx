import {Step} from 'react-joyride';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';

export enum Steps {
  AddTestSpec = 'add-test-spec',
  Graph = 'graph',
  Switcher = 'switcher',
  MoreData = 'more-data',
  SpanDetails = 'assertions',
  RunButton = 'run-button',
  MetaDetails = 'meta-details',
}

export const switchTraceMode = (index: number) => {
  const elementNodeListOfElement = (document.querySelectorAll('.ant-tabs-tab') as NodeListOf<HTMLElement>)[index];
  if (elementNodeListOfElement !== null) {
    elementNodeListOfElement.click();
  }
};

const StepList: Step[] = [
  {
    title: 'Response',
    target: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.Graph),
    content: 'View various responses here. When the test is finished, you will get the following results.',
    placement: 'left',
    disableBeacon: true,
  },
  {
    title: 'Add more data',
    target: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.MoreData),
    placement: 'right',
    content:
      'You can change the trigger by altering the details and saving. This will rerun the test with the updated trigger.',
    disableBeacon: true,
  },
  {
    title: 'Mode Switcher',
    target: '.ant-tabs-nav-wrap',
    content: 'Click on the Trace tab to open the Trace Details screen or Test tab for adding test specifications.',
    disableBeacon: true,
  },
  {
    title: 'Span Details',
    target: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.SpanDetails),
    content:
      'Click on the panel to see details about the selected span. These span attributes are grouped into tabs based on the type of span. ',
    disableBeacon: true,
    placement: 'right',
  },
  {
    title: 'Test Spec',
    target: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.AddTestSpec),
    content:
      'Test Specifications can be added to set assertions to run against one or more spans in the trace. If test specs have already been added to a test, there will be a list on the Test screen.',
    placement: 'left',
    disableBeacon: true,
  },
];

export default StepList;
