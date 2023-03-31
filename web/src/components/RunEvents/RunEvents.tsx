import {TestRunStage} from 'constants/TestRunEvents.constants';
import TestRunEvent from 'models/TestRunEvent.model';
import TestRunService from 'services/TestRun.service';
import {TTestRunState} from 'types/TestRun.types';
import RunEventsTrace from './RunEventsTrace';
import RunEventsTrigger from './RunEventsTrigger';

export interface IPropsComponent {
  events: TestRunEvent[];
  state: TTestRunState;
}

interface IProps extends IPropsComponent {
  stage: TestRunStage;
}

const ComponentMap: Record<TestRunStage, (props: IPropsComponent) => React.ReactElement> = {
  [TestRunStage.Trigger]: RunEventsTrigger,
  [TestRunStage.Trace]: RunEventsTrace,
  [TestRunStage.Test]: RunEventsTrace,
};

const RunEvents = ({events, stage, state}: IProps) => {
  const Component = ComponentMap[stage];
  const filteredEvents = TestRunService.getTestRunEventsByStage(events, stage);

  return <Component events={filteredEvents} state={state} />;
};

export default RunEvents;
