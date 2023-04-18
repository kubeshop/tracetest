import {TestState} from 'constants/TestRun.constants';
import {TraceEventType} from 'constants/TestRunEvents.constants';
import TestRunService from 'services/TestRun.service';
import RunEvent, {IPropsEvent} from './RunEvent';
import RunEventDataStore from './RunEventDataStore';
import RunEventPolling from './RunEventPolling';
import {IPropsComponent} from './RunEvents';
import * as S from './RunEvents.styled';
import FailedTraceHeader from './TraceHeader/FailedTraceHeader';
import FailedTriggerHeader from './TraceHeader/FailedTriggerHeader';
import LoadingHeader from './TraceHeader/LoadingHeader';
import StoppedHeader from './TraceHeader/StoppedHeader';

type TestStateType = TestState.TRIGGER_FAILED | TestState.TRACE_FAILED | TestState.STOPPED;

const HeaderComponentMap: Record<TestStateType, () => React.ReactElement> = {
  [TestState.TRIGGER_FAILED]: FailedTriggerHeader,
  [TestState.TRACE_FAILED]: FailedTraceHeader,
  [TestState.STOPPED]: StoppedHeader,
};

export type TraceEventTypeWithoutFetching = Exclude<
  TraceEventType,
  TraceEventType.FETCHING_START | TraceEventType.FETCHING_ERROR | TraceEventType.FETCHING_SUCCESS
>;

const ComponentMap: Record<TraceEventTypeWithoutFetching, (props: IPropsEvent) => React.ReactElement> = {
  [TraceEventType.DATA_STORE_CONNECTION_INFO]: RunEventDataStore,
  [TraceEventType.POLLING_ITERATION_INFO]: RunEventPolling,
};

const RunEventsTrace = ({events, state}: IPropsComponent) => {
  const filteredEvents = TestRunService.getTestRunTraceEvents(events);
  const HeaderComponent = HeaderComponentMap[state as TestStateType] ?? LoadingHeader;

  return (
    <S.Container $hasScroll>
      <HeaderComponent />

      <S.ListContainer>
        {filteredEvents.map(event => {
          const Component = ComponentMap[event.type as TraceEventTypeWithoutFetching] ?? RunEvent;
          return <Component event={event} key={event.type} />;
        })}
      </S.ListContainer>
    </S.Container>
  );
};

export default RunEventsTrace;
