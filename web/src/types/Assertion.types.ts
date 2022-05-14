import {CompareOperator} from 'constants/Operator.constants';
import {LOCATION_NAME} from 'constants/Span.constants';
import {SpanAttributeType} from 'constants/SpanAttribute.constants';
import {ISpan} from './Span.types';

export interface IAssertion {
  assertionId: string;
  selectors: Array<IItemSelector>;
  spanAssertions: Array<ISpanSelector>;
}

export interface IItemSelector {
  locationName: LOCATION_NAME;
  propertyName: string;
  value: string;
  valueType: string;
}

export interface ISpanSelector {
  spanAssertionId?: string;
  locationName: LOCATION_NAME;
  propertyName: string;
  valueType: SpanAttributeType;
  operator: CompareOperator;
  comparisonValue: string;
}

export type TAssertionResultList = Array<{
  assertionId: string;
  spanAssertionResults: ISpanAssertionResult2[];
}>;

export interface ITestAssertionResult {
  assertionResultState: boolean;
  assertionResult: TAssertionResultList;
}

export interface IAssertionResult {
  spanListAssertionResult: {
    span: ISpan;
    resultList: ISpanAssertionResult[];
  }[];
  assertion: IAssertion;
}

export interface ISpanAssertionResult extends ISpanSelector {
  hasPassed: boolean;
  actualValue: string;
  spanId: string;
}

export interface ISpanAssertionResult2 {
  spanAssertionId: string;
  spanId: string;
  passed: boolean;
  observedValue: string;
}

export interface IAssertionResultList {
  assertion: IAssertion;
  assertionResultList: ISpanAssertionResult[];
}
