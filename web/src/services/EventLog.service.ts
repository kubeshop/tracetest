import {parseISO, formatISO} from 'date-fns';
import TestRunEvent, {PollingInfo} from 'models/TestRunEvent.model';
import ConnectionResult from 'models/ConnectionResult.model';
import {TraceEventType} from 'constants/TestRunEvents.constants';

type TEventToStringFn = (event: TestRunEvent) => string;

const eventToString = ({title, description}: TestRunEvent): string => {
  return `${title} - ${description}`;
};

const dataStoreEventToString = (event: TestRunEvent): string => {
  const {dataStoreConnection: {allPassed, ...dataStoreConnection} = ConnectionResult({})} = event;
  const baseText = eventToString(event);
  const configValidText = allPassed ? 'Data store configuration is valid.' : 'Data store configuration is not valid.';

  const connectionStepsDetailsText = Object.entries(dataStoreConnection || {})
    .map(([key, {message, error}]) => `${key.toUpperCase()} - ${message} ${error ? ` - ${error}` : ''}`, '')
    .join(' - ');

  return `${baseText} - ${configValidText} - ${connectionStepsDetailsText}`;
};

const pollingEventToString = (event: TestRunEvent): string => {
  const {polling: {type: pollingType, isComplete, periodic} = PollingInfo({})} = event;
  const baseText = eventToString(event);
  const pollingTypeText = `Polling type: ${pollingType}`;
  const pollingCompleteText = `Polling complete: ${isComplete}`;
  const periodicText = `Periodic polling - number of spans: ${periodic?.numberSpans}, number of iterations: ${periodic?.numberIterations}`;

  return `${baseText} - ${pollingTypeText} - ${pollingCompleteText} - ${periodicText}`;
};

const eventToStringMap: Record<string, TEventToStringFn> = {
  [TraceEventType.DATA_STORE_CONNECTION_INFO]: dataStoreEventToString,
  [TraceEventType.POLLING_ITERATION_INFO]: pollingEventToString,
};

const EventLogService = () => ({
  detailsToString(event: TestRunEvent) {
    const eventToStringFn = eventToStringMap[event.type] || eventToString;

    return eventToStringFn(event);
  },
  typeToString({type, createdAt}: TestRunEvent) {
    const createdAtDate = parseISO(createdAt);

    return `[${formatISO(createdAtDate)} - ${type}]`;
  },
  listToString(events: TestRunEvent[]) {
    return events.map(event => `${this.typeToString(event)} ${this.detailsToString(event)}`).join('\r');
  },
});

export default EventLogService();
