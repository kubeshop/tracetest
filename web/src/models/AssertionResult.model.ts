import {TAssertionResult, TRawAssertionResult} from '../types/Assertion.types';
import Assertion from './Assertion.model';
import AssertionSpanResult from './AssertionSpanResult.model';

const AssertionResult = ({allPassed = false, assertion, spanResults = []}: TRawAssertionResult): TAssertionResult => {
  return {
    allPassed,
    assertion: Assertion(assertion!),
    spanResults: spanResults.map(spanResult => AssertionSpanResult(spanResult)),
  };
};

export default AssertionResult;
