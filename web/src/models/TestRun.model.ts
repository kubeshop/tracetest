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
  request,
  response,
  testVersion = 1,
  exectutionTime = 0,
  obtainedTraceAt = '',
  serviceTriggerCompletedAt = '',
  serviceTriggeredAt = '',
}: TRawTestRun): TTestRun => {
  return {
    obtainedTraceAt,
    serviceTriggerCompletedAt,
    serviceTriggeredAt,
    executionTime: exectutionTime,
    exectutionTime,
    lastErrorState,
    request,
    response,
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
  };
};

export default TestRun;
