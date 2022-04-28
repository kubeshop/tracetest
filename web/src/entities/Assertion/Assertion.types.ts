import {CompareOperator} from '../Operator/Operator.constants';
import {LOCATION_NAME} from '../Span/Span.constants';
import {TResourceSpan} from '../Span/Span.types';
import {SpanAttributeType} from '../SpanAttribute/SpanAttribute.constants';

export type TAssertion = {
  assertionId: string;
  selectors: Array<TItemSelector>;
  spanAssertions: Array<TSpanSelector>;
};

export type TItemSelector = {
  locationName: LOCATION_NAME;
  propertyName: string;
  value: string;
  valueType: string;
};

export type TSpanSelector = {
  spanAssertionId?: string;
  locationName: LOCATION_NAME;
  propertyName: string;
  valueType: SpanAttributeType;
  operator: CompareOperator;
  comparisonValue: string;
};

export type TAssertionResultList = Array<{
  assertionId: string;
  spanAssertionResults: TSpanAssertionResult2[];
}>;

export type TTestAssertionResult = {
  assertionResultState: boolean;
  assertionResult: TAssertionResultList;
};

export type TAssertionResult = {
  spanListAssertionResult: {
    span: TResourceSpan;
    resultList: TSpanAssertionResult[];
  }[];
  assertion: TAssertion;
};

export type TSpanAssertionResult = TSpanSelector & {
  hasPassed: boolean;
  actualValue: string;
  spanId: string;
};

export type TSpanAssertionResult2 = {
  spanAssertionId: string;
  spanId: string;
  passed: boolean;
  observedValue: string;
};
