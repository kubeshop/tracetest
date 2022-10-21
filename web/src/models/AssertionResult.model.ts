import {TAssertionResult, TRawAssertionResult} from '../types/Assertion.types';
import AssertionSpanResult from './AssertionSpanResult.model';

const AssertionResult = ({allPassed = false, assertion = '', spanResults = []}: TRawAssertionResult): TAssertionResult => {
  return {
    allPassed,
    assertion,
    spanResults: spanResults.map(spanResult => AssertionSpanResult(spanResult)),
  };
};

export default AssertionResult;
