import {PseudoSelector} from 'constants/Operator.constants';
import {Model, TTestSchemas} from './Common.types';
import {TCompareOperatorSymbol} from './Operator.types';
import {TSpanFlatAttribute} from './Span.types';

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

export type TStructuredAssertion = {
  left: string;
  comparator: TCompareOperatorSymbol;
  right: string;
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

export type TResultAssertionsSummary = {
  failed: IResult[];
  passed: IResult[];
};

export type TResultAssertions = Record<string, TResultAssertionsSummary>;

export interface ICheckResult {
  result: TAssertionSpanResult;
  assertion: string;
}
