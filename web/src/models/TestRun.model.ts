import {TRawAssertionResults} from '../types/Assertion.types';
import {TRawTestRun, TTestRun} from '../types/TestRun.types';
import AssertionResults from './AssertionResults.model';
import Trace from './Trace.model';

const getTestResultCount = (
  {results: resultList = []}: TRawAssertionResults = {},
  type: 'all' | 'passed' | 'failed' = 'all'
) => {
  const spanAssertionList = resultList.flatMap(({results = []}) =>
    results.flatMap(({spanResults = []}) => spanResults)
  );

  if (type === 'all') return spanAssertionList.length;

  return spanAssertionList.filter(({passed}) => {
    switch (type) {
      case 'failed': {
        return !passed;
      }

      case 'passed':
      default: {
        return passed;
      }
    }
  }).length;
};

const TestRun = ({
  id = '',
  traceId = '',
  spanId = '',
  state = 'CREATED',
  createdAt = '',
  completedAt = '',
  trace,
  result,
  lastErrorState,
  trigger,
  triggerResult,
  testVersion = 1,
  executionTime = 0,
  obtainedTraceAt = '',
  serviceTriggerCompletedAt = '',
  serviceTriggeredAt = '',
  metadata = {},
}: TRawTestRun): TTestRun => {
  return {
    obtainedTraceAt,
    serviceTriggerCompletedAt,
    serviceTriggeredAt,
    executionTime,
    lastErrorState,
    trigger,
    triggerResult,
    createdAt,
    completedAt,
    result: AssertionResults(result!),
    id,
    traceId,
    spanId,
    state,
    testVersion,
    trace: trace ? Trace(trace) : undefined,
    totalAssertionCount: getTestResultCount(result),
    failedAssertionCount: getTestResultCount(result, 'failed'),
    passedAssertionCount: getTestResultCount(result, 'passed'),
    metadata,
  };
};

export default TestRun;
