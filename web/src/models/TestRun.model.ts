import {differenceInSeconds} from 'date-fns';
import {TRawAssertionResults} from '../types/Assertion.types';
import {TRawTestRun, TTestRun} from '../types/TestRun.types';
import AssertionResults from './AssertionResults.model';
import Trace from './Trace.model';

const getExecutionTime = (createdAt?: string, completedAt?: string) => {
  if (!createdAt || !completedAt) return 0;

  return differenceInSeconds(new Date(completedAt), new Date(createdAt)) + 1;
};

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
}: TRawTestRun): TTestRun => {
  return {
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
    trace: trace ? Trace(trace) : undefined,
    totalAssertionCount: getTestResultCount(result),
    failedAssertionCount: getTestResultCount(result, 'failed'),
    passedAssertionCount: getTestResultCount(result, 'passed'),
    executionTime: getExecutionTime(createdAt, completedAt),
  };
};

export default TestRun;
