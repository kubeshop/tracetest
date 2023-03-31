import {TestRunStage, TraceEventType} from 'constants/TestRunEvents.constants';
import {filter, findLastIndex} from 'lodash';
import TestRunEvent from 'models/TestRunEvent.model';

const TestRunService = () => ({
  getTestRunEventsByStage(events: TestRunEvent[], stage: TestRunStage) {
    return events.filter(event => event.stage === stage);
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
