import TestRunEvent, {TestRunStage} from 'models/TestRunEvent.model';
import TestRunService from 'services/TestRun.service';
import RunEventsTrigger from './RunEventsTrigger';

export interface IPropsComponent {
  events: TestRunEvent[];
}

interface IProps extends IPropsComponent {
  stage: TestRunStage;
}

const ComponentMap: Record<TestRunStage, (props: IPropsComponent) => React.ReactElement> = {
  [TestRunStage.Trigger]: RunEventsTrigger,
  [TestRunStage.Trace]: RunEventsTrigger,
  [TestRunStage.Test]: RunEventsTrigger,
};

const RunEvents = ({events, stage}: IProps) => {
  const Component = ComponentMap[stage];
  const filteredEvents = TestRunService.getTestRunEventsByStage(events, stage);

  return <Component events={filteredEvents} />;
};

export default RunEvents;
