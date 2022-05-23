import {TAssertionSpanResult, TRawAssertionSpanResult} from '../types/Assertion.types';

const AssertionSpanResult = ({
  spanId = '',
  observedValue = '',
  passed = false,
  error = '',
}: TRawAssertionSpanResult): TAssertionSpanResult => {
  return {
    spanId,
    observedValue,
    passed,
    error,
  };
};

export default AssertionSpanResult;
