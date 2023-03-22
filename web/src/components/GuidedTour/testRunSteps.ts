import {Step} from 'react-joyride';
import GuidedTourService from 'services/GuidedTour.service';

export enum StepsID {
  Response = 'testRun_response',
  Trigger = 'testRun_trigger',
  SpanDetails = 'testRun_span-details',
  TestSpecs = 'testRun_test-specs',
}

const Steps: Step[] = [
  {
    title: 'Trigger',
    target: GuidedTourService.getStepSelector(StepsID.Trigger),
    placement: 'right',
    content:
      'You can change the trigger by altering the details and saving. This will rerun the test with the updated trigger.',
    disableBeacon: true,
  },
  {
    title: 'Response',
    target: GuidedTourService.getStepSelector(StepsID.Response),
    content: 'View various responses here. When the test is finished, you will get the following results.',
    placement: 'left',
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
    target: GuidedTourService.getStepSelector(StepsID.SpanDetails),
    content: 'Click on the panel to see details about the selected span. Here you can find and search span attributes.',
    disableBeacon: true,
    placement: 'right',
  },
  {
    title: 'Test Specs',
    target: GuidedTourService.getStepSelector(StepsID.TestSpecs),
    content:
      'Test Specifications can be added to set assertions to run against one or more spans in the trace. If test specs have already been added to a test, there will be a list on the Test screen.',
    placement: 'left',
    disableBeacon: true,
  },
];

export const switchTestRunMode = (index: number) => {
  const elementNodeListOfElement = (document.querySelectorAll('.ant-tabs-tab div a') as NodeListOf<HTMLElement>)[index];
  if (elementNodeListOfElement !== null) {
    elementNodeListOfElement.click();
  }
};

export default Steps;
