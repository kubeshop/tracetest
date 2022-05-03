import {Schemas} from './Common.types';
import {TSpan} from './Span.types';

export type TAssertion = Schemas['Assertion'];
export type TItemSelector = Schemas['SelectorItem'];
export type TSpanSelector = Schemas['SpanAssertion'];

export type TAssertionResultList = Array<Schemas['AssertionResult']>;

export interface TTestAssertionResult {
  assertionResultState: boolean;
  assertionResult: TAssertionResultList;
}

export type TAssertionResult = {
  spanListAssertionResult: {
    span: TSpan;
    resultList: TSpanAssertionResult[];
  }[];
  assertion: TAssertion;
};

export interface TSpanAssertionResult extends TSpanSelector {
  hasPassed: boolean;
  actualValue: string;
  spanId: string;
}

export type TSpanAssertionResult2 = Schemas['SpanAssertionResult'];
