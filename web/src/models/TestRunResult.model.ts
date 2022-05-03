import {TRawTestRunResult, TTestRunResult} from '../types/TestRunResult.types';
import Trace from './Trace.model';

const TestRunResult = (rawTestRunResult: TRawTestRunResult): TTestRunResult => {
  return {
    ...rawTestRunResult,
    trace: rawTestRunResult.trace ? Trace(rawTestRunResult.trace) : undefined,
  };
};

export default TestRunResult;
