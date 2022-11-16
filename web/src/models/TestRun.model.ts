import {TRawAssertionResults} from 'types/Assertion.types';
import {TRawTestRun, TTestRun} from 'types/TestRun.types';
import AssertionResults from './AssertionResults.model';
import {TestRunOutput} from './TestOutput.model';
import Trace from './Trace.model';
import TriggerResult from './TriggerResult.model';

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
  triggerResult: rawTriggerResult,
  testVersion = 1,
  executionTime = 0,
  triggerTime = 0,
  obtainedTraceAt = '',
  serviceTriggerCompletedAt = '',
  serviceTriggeredAt = '',
  metadata = {},
  outputs = [],
}: TRawTestRun): TTestRun => {
  return {
    obtainedTraceAt,
    serviceTriggerCompletedAt,
    serviceTriggeredAt,
    executionTime,
    triggerTime,
    lastErrorState,
    triggerResult: rawTriggerResult ? TriggerResult(rawTriggerResult) : undefined,
    createdAt,
    completedAt,
    result: AssertionResults(result || {}),
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
    outputs: outputs?.map(rawOuput => TestRunOutput(rawOuput)),
  };
};

export default TestRun;
