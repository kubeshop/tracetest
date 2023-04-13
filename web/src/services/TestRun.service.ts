import {TestRunStage, TraceEventType} from 'constants/TestRunEvents.constants';
import {filter, findLastIndex, flow} from 'lodash';
import {isRunStateStopped, isRunStateSucceeded} from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import {TTestRunState} from 'types/TestRun.types';

const TestRunService = () => ({
  shouldDisplayTraceEvents(state: TTestRunState, numberOfSpans: number) {
    const isStateSucceededOrStopped = isRunStateSucceeded(state) || isRunStateStopped(state);
    return !isStateSucceededOrStopped || (!numberOfSpans && isStateSucceededOrStopped);
  },

  getTestRunEventsByStage(events: TestRunEvent[], stage: TestRunStage) {
    return events.filter(event => event.stage === stage);
  },

  getTestRunTraceEvents(events: TestRunEvent[]): TestRunEvent[] {
    return flow([this.getTestRunEventsWithoutFetching, this.getTestRunEventsWithLastPolling])(events);
  },

  getTestRunEventsWithoutFetching(events: TestRunEvent[]) {
    return filter(
      events,
      event =>
        !(
          [TraceEventType.FETCHING_START, TraceEventType.FETCHING_ERROR, TraceEventType.FETCHING_SUCCESS] as string[]
        ).includes(event.type)
    );
  },

  getTestRunEventsWithLastPolling(events: TestRunEvent[]) {
    const lastPollingIndex = findLastIndex(events, event => event.type === TraceEventType.POLLING_ITERATION_INFO);
    if (lastPollingIndex === -1) return events;

    const eventsWithoutPolling = filter(events, event => event.type !== TraceEventType.POLLING_ITERATION_INFO);
    const newIndex = lastPollingIndex - (events.length - eventsWithoutPolling.length) + 1;
    eventsWithoutPolling.splice(newIndex, 0, events[lastPollingIndex]);

    return eventsWithoutPolling;
  },
});

export default TestRunService();
