import {differenceInSeconds} from 'date-fns';
import {TAssertionResultList} from '../types/Assertion.types';
import {IRawTestRunResult, ITestRunResult} from '../types/TestRunResult.types';
import Trace from './Trace.model';

const getExecutionTime = (createdAt?: string, completedAt?: string) => {
  if (!createdAt || !completedAt) return 0;

  return differenceInSeconds(new Date(completedAt), new Date(createdAt)) + 1;
};

const getTestResultCount = (assertionResultList: TAssertionResultList = [], type: 'all' | 'passed' | 'failed' = 'all') => {
  const spanAssertionList = assertionResultList.flatMap(({spanAssertionResults}) => spanAssertionResults);

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

const TestRunResult = (rawTestRunResult: IRawTestRunResult): ITestRunResult => {
  return {
    ...rawTestRunResult,
    trace: rawTestRunResult.trace ? Trace(rawTestRunResult.trace) : undefined,
    totalAssertionCount: getTestResultCount(rawTestRunResult.assertionResult, 'all'),
    failedAssertionCount: getTestResultCount(rawTestRunResult.assertionResult, 'failed'),
    passedAssertionCount: getTestResultCount(rawTestRunResult.assertionResult, 'passed'),
    executionTime: getExecutionTime(rawTestRunResult.createdAt, rawTestRunResult.completedAt),
  };
};

export default TestRunResult;
