import {StepType} from '@reactour/tour';
import {Typography} from 'antd';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {StepContent} from './StepContent';

export enum Steps {
  Diagram = 'diagram',
  SpanDetail = 'spanDetail',
  TestResults = 'testResults',
  Timeline = 'timeline',
  Graph = 'graph',
  Details = 'details',
  Switcher = 'switcher',
  Assertions = 'assertions',
}

const StepList: StepType[] = [
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.Graph),
    content: ({setIsOpen}) => (
      <StepContent title="Test View" setIsOpen={setIsOpen}>
        <Typography.Text>
          The trace view window shows you the trace generated from the triggering transaction. This is the Graph view
          showing the calling relationship between services.
        </Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.Switcher),
    content: ({setIsOpen}) => (
      <StepContent title="Switcher" setIsOpen={setIsOpen}>
        <Typography.Text>
          You can view the complete trace in a graph view or you can see the trace in a timeline view by clicking the
          switcher
        </Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.Details),
    content: ({setIsOpen}) => (
      <StepContent title="Span Details" setIsOpen={setIsOpen}>
        <Typography.Text>
          Details about the selected span are shown here. They are grouped into tabs based on the type of span. The
          {` 'Attribute list' `} shows all of the attributes
        </Typography.Text>
      </StepContent>
    ),
  },
  {
    selector: GuidedTourService.getSelectorStep(GuidedTours.Trace, Steps.Assertions),
    content: ({setIsOpen}) => (
      <StepContent title="Adding Assertions" setIsOpen={setIsOpen}>
        <Typography.Text>
          You can add an assertion to the attribute on any span by hovering over it and click the plus sign (+) or click
          Add Assertion button below. Assertions may be added to a trace to set a value for a step in the trace to
          determine success or failure.
        </Typography.Text>
      </StepContent>
    ),
  },
];

export default StepList;
