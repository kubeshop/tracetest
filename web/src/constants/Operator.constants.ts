import {TCompareOperatorName, TCompareOperatorSymbol} from '../types/Operator.types';

export enum CompareOperator {
  EQUALS = 'EQUALS',
  NOTEQUALS = 'NOTEQUALS',
  LESSTHAN = 'LESSTHAN',
  GREATERTHAN = 'GREATERTHAN',
  GREATOREQUALS = 'GREATOREQUALS',
  LESSOREQUAL = 'LESSOREQUAL',
  CONTAINS = 'CONTAINS',
  NOTCONTAINS = 'NOTCONTAINS',
}

export const CompareOperatorNameMap: Record<CompareOperator, TCompareOperatorName> = {
  [CompareOperator.EQUALS]: 'equals',
  [CompareOperator.NOTEQUALS]: 'not equals',
  [CompareOperator.LESSTHAN]: 'less than',
  [CompareOperator.GREATERTHAN]: 'greater than',
  [CompareOperator.GREATOREQUALS]: 'greater or equals',
  [CompareOperator.LESSOREQUAL]: 'less or equals',
  [CompareOperator.CONTAINS]: 'contains',
  [CompareOperator.NOTCONTAINS]: 'does not contain',
};

export const CompareOperatorSymbolMap: Record<CompareOperator, TCompareOperatorSymbol> = {
  [CompareOperator.EQUALS]: '=',
  [CompareOperator.LESSTHAN]: '<',
  [CompareOperator.GREATERTHAN]: '>',
  [CompareOperator.NOTEQUALS]: '!=',
  [CompareOperator.GREATOREQUALS]: '>=',
  [CompareOperator.LESSOREQUAL]: '<=',
  [CompareOperator.CONTAINS]: 'contains',
  [CompareOperator.NOTCONTAINS]: 'not-contains',
};

export const CompareOperatorSymbolNameMap: Record<TCompareOperatorSymbol, TCompareOperatorName> = {
  '=': 'equals',
  '!=': 'not equals',
  '<': 'less than',
  '>': 'greater than',
  '>=': 'greater or equals',
  '<=': 'less or equals',
  contains: 'contains',
  'not-contains': 'does not contain',
};

export enum PseudoSelector {
  FIRST = ':first',
  LAST = ':last',
  NTH = ':nth_child',
  ALL = '',
}

export const OperatorRegexp = /= | != | < | > | >= | <= | contains | not-contains/gi;
