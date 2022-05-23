import {CompareOperator} from 'constants/Operator.constants';
import {Model, TTestSchemas} from './Common.types';
import {TSpanFlatAttribute} from './Span.types';

export type TRawAssertion = TTestSchemas['Assertion'];

export type TAssertion = Model<TRawAssertion, {}>;

export type TSpanSelector = Model<
  TSpanFlatAttribute,
  {
    operator: CompareOperator;
  }
>;

export type TRawAssertionResults = TTestSchemas['AssertionResults'];

export type TAssertionResultEntry = {
  id: string;
  selector: string;
  resultList: TAssertionResult[];
};

export type TAssertionResults = Model<
  TRawAssertionResults,
  {
    allPassed: boolean;
    results?: never;
    resultList: TAssertionResultEntry[];
  }
>;

export type TRawAssertionResult = TTestSchemas['AssertionResult'];
export type TAssertionResult = Model<
  TRawAssertionResult,
  {
    assertion: TAssertion;
    spanResults: TAssertionSpanResult[];
  }
>;

export type TRawAssertionSpanResult = TTestSchemas['AssertionSpanResult'];
export type TAssertionSpanResult = Model<TRawAssertionSpanResult, {}>;
