import {PseudoSelector} from '../constants/Operator.constants';
import {Model, TTestSchemas} from './Common.types';
import {TCompareOperatorSymbol} from './Operator.types';
import {TSpanFlatAttribute} from './Span.types';

export type TRawAssertion = TTestSchemas['Assertion'];

export type TAssertion = Model<TRawAssertion, {}>;

export type TSpanSelector = Model<
  TSpanFlatAttribute,
  {
    operator: TCompareOperatorSymbol;
  }
>;

export type TPseudoSelector = {
  selector: PseudoSelector;
  number?: number;
};

export type TRawAssertionResults = TTestSchemas['AssertionResults'];

export type TAssertionResultEntry = {
  id: string;
  selector: string;
  originalSelector: string;
  spanIds: string[];
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

export interface IResult {
  id: string;
  label: string;
  assertionResult: TAssertionResultEntry;
}

export type TResultAssertions = Record<
  string,
  {
    failed: IResult[];
    passed: IResult[];
  }
>;

export interface ICheckResult {
  result: TAssertionSpanResult;
  assertion: TAssertion;
}
