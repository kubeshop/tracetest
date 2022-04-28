import {TCompareOperatorName, TCompareOperatorSymbol} from './Operator.types';

export const enum CompareOperator {
  EQUALS = 'EQUALS',
  NOTEQUALS = 'NOTEQUALS',
  LESSTHAN = 'LESSTHAN',
  GREATERTHAN = 'GREATERTHAN',
  GREATOREQUALS = 'GREATOREQUALS',
  LESSOREQUAL = 'LESSOREQUAL',
  CONTAINS = 'CONTAINS',
}

export const CompareOperatorNameMap: Record<CompareOperator, TCompareOperatorName> = {
  [CompareOperator.EQUALS]: 'eq',
  [CompareOperator.NOTEQUALS]: 'ne',
  [CompareOperator.LESSTHAN]: 'lt',
  [CompareOperator.GREATERTHAN]: 'gt',
  [CompareOperator.GREATOREQUALS]: 'gte',
  [CompareOperator.LESSOREQUAL]: 'lte',
  [CompareOperator.CONTAINS]: 'contains',
};

export const CompareOperatorSymbolMap: Record<CompareOperator, TCompareOperatorSymbol> = {
  [CompareOperator.EQUALS]: '==',
  [CompareOperator.LESSTHAN]: '<',
  [CompareOperator.GREATERTHAN]: '>',
  [CompareOperator.NOTEQUALS]: '!=',
  [CompareOperator.GREATOREQUALS]: '>=',
  [CompareOperator.LESSOREQUAL]: '<=',
  [CompareOperator.CONTAINS]: 'contains',
};
