import {Model, TTestSchemas} from 'types/Common.types';

type AssertionSpanResult = Model<TRawAssertionSpanResult, {}>;
export type TRawAssertionSpanResult = TTestSchemas['AssertionSpanResult'];

const AssertionSpanResult = ({
  spanId = '',
  observedValue = '',
  passed = false,
  error = '',
}: TRawAssertionSpanResult): AssertionSpanResult => {
  return {
    spanId,
    observedValue,
    passed,
    error,
  };
};

export default AssertionSpanResult;
