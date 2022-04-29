import {IRawTestRunResult, ITestRunResult} from '../types/TestRunResult.types';
import Trace from './Trace.model';

const TestRunResult = (rawTestRunResult: IRawTestRunResult): ITestRunResult => {
  return {
    ...rawTestRunResult,
    trace: rawTestRunResult.trace ? Trace(rawTestRunResult.trace) : undefined,
  };
};

export default TestRunResult;
