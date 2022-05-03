import { Schemas } from '../types/Common.types';
import {TCompareOperatorName, TCompareOperatorSymbol} from '../types/Operator.types';

export const CompareOperatorNameMap: Record<Schemas['SpanAssertion']['operator'], TCompareOperatorName> = {
  EQUALS: 'eq',
  NOTEQUALS: 'ne',
  LESSTHAN: 'lt',
  GREATERTHAN: 'gt',
  GREATOREQUALS: 'gte',
  LESSOREQUAL: 'lte',
  CONTAINS: 'contains',
};

export const CompareOperatorSymbolMap: Record<Schemas['SpanAssertion']['operator'], TCompareOperatorSymbol> = {
  EQUALS: '==',
  LESSTHAN: '<',
  GREATERTHAN: '>',
  NOTEQUALS: '!=',
  GREATOREQUALS: '>=',
  LESSOREQUAL: '<=',
  CONTAINS: 'contains',
};
