import {TestState} from 'constants/TestRun.constants';
import {Model, Modify, TTestSchemas, TTriggerSchemas} from 'types/Common.types';
import {TTestRunState} from 'types/TestRun.types';
import AssertionResults, {TRawAssertionResults} from './AssertionResults.model';
import VariableSet from './VariableSet.model';
import LinterResult from './LinterResult.model';
import TestRunOutput from './TestRunOutput.model';
import Trace from './Trace.model';
import TriggerResult from './TriggerResult.model';
import RequiredGatesResult from './RequiredGatesResult.model';
import RunMetadata from './RunMetadata.model';

export type TRawTestRun = Modify<
  TTestSchemas['TestRun'],
  {
    state?: TTestRunState;
  }
>;
type TestRun = Model<
  TRawTestRun,
  {
    result: AssertionResults;
    trace: Trace;
    totalAssertionCount: number;
    failedAssertionCount: number;
    passedAssertionCount: number;
    executionTime: number;
    triggerTime: number;
    lastErrorState?: string;
    trigger?: TTriggerSchemas['Trigger'];
    resolvedTrigger?: TTriggerSchemas['Trigger'];
    triggerResult?: TriggerResult;
    outputs?: TestRunOutput[];
    variableSet?: VariableSet;
    state: TTestRunState;
    linter: LinterResult;
    requiredGatesResult: RequiredGatesResult;
    testSuiteRunId: number;
    metadata: RunMetadata;
  }
>;

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

export function isRunStateFinished(state: TTestRunState) {
  return (
    [
      TestState.FINISHED,
      TestState.STOPPED,
      TestState.TRIGGER_FAILED,
      TestState.TRACE_FAILED,
      TestState.ASSERTION_FAILED,
      TestState.ANALYZING_ERROR,
    ] as string[]
  ).includes(state);
}

export function isRunStateFailed(state: TTestRunState) {
  return ([TestState.TRIGGER_FAILED, TestState.TRACE_FAILED, TestState.ASSERTION_FAILED] as string[]).includes(state);
}

export function isRunStateSucceeded(state: TTestRunState) {
  return state === TestState.FINISHED;
}

export function isRunStateStopped(state: TTestRunState) {
  return state === TestState.STOPPED;
}

export function isRunStateAnalyzingError(state: TTestRunState) {
  return state === TestState.ANALYZING_ERROR;
}

export function isRunPollingState(state: TTestRunState) {
  return state === TestState.AWAITING_TRACE;
}

const TestRun = ({
  id = 0,
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
  variableSet = {},
  testSuiteId = '',
  testSuiteRunId = 0,
  linter = {},
  requiredGatesResult = {required: [], failed: [], passed: true},
}: TRawTestRun): TestRun => {
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
    trace: trace ? Trace(trace) : Trace(),
    totalAssertionCount: getTestResultCount(result),
    failedAssertionCount: getTestResultCount(result, 'failed'),
    passedAssertionCount: getTestResultCount(result, 'passed'),
    metadata: RunMetadata(metadata),
    outputs: outputs?.map(rawOutput => TestRunOutput(rawOutput)),
    variableSet: VariableSet.fromRun(variableSet),
    testSuiteId,
    testSuiteRunId: testSuiteRunId ? Number(testSuiteRunId) : 0,
    linter: LinterResult(linter),
    requiredGatesResult: RequiredGatesResult(requiredGatesResult),
  };
};

export default TestRun;
