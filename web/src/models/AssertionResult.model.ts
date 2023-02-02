import {Model, TTestSchemas} from 'types/Common.types';
import AssertionSpanResult from './AssertionSpanResult.model';

export type TRawAssertionResult = TTestSchemas['AssertionResult'];
type AssertionResult = Model<
  TRawAssertionResult,
  {
    spanResults: AssertionSpanResult[];
  }
>;

const AssertionResult = ({
  allPassed = false,
  assertion = '',
  spanResults = [],
}: TRawAssertionResult): AssertionResult => {
  return {
    allPassed,
    assertion,
    spanResults: spanResults.map(spanResult => AssertionSpanResult(spanResult)),
  };
};

export default AssertionResult;
