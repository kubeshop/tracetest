import {CompareOperatorSymbolMap, PseudoSelector} from './Operator.constants';

export enum Tokens {
  Span = 'Span',
  SpanMatch = 'SpanMatch',
  BaseComparator = 'BaseComparator',
  ComparatorValue = 'ComparatorValue',
  Identifier = 'Identifier',
  Operator = 'Operator',
  Number = 'Number',
  String = 'String',
  ClosingBracket = 'ClosingBracket',
  Program = 'Program',
}

export const completeIdentifierAfter: string[] = [
  Tokens.Span,
  Tokens.SpanMatch,
  Tokens.BaseComparator,
  Tokens.ComparatorValue,
  Tokens.Identifier,
];
export const completeOperatorAfter: string[] = [Tokens.Identifier];
export const completeValueAfter: string[] = [Tokens.ComparatorValue];

export const completePseudoSelectorAfter: string[] = [Tokens.ClosingBracket];

export const operatorList = [
  {
    label: CompareOperatorSymbolMap.EQUALS,
    type: 'operatorKeyword',
  },
  {
    label: CompareOperatorSymbolMap.CONTAINS,
    type: 'operatorKeyword',
  },
];

export const pseudoSelectorList = [
  {
    label: PseudoSelector.FIRST,
    type: 'operatorKeyword',
  },
  {
    label: PseudoSelector.LAST,
    type: 'operatorKeyword',
  },
  {
    label: PseudoSelector.NTH,
    type: 'operatorKeyword',
  },
];
