import {CompareOperatorSymbolMap, PseudoSelector} from './Operator.constants';

export enum Tokens {
  Span = 'Span',
  SpanMatch = 'SpanMatch',
  BaseComparator = 'BaseComparator',
  ComparatorValue = 'ComparatorValue',
  ComposedValue = 'ComposedValue',
  Identifier = 'Identifier',
  Comparator = 'Comparator',
  Operator = 'Operator',
  Number = 'Number',
  String = 'String',
  ClosingBracket = 'ClosingBracket',
  Program = 'Program',
  Source = 'Source',
  OutsideInput = 'OutsideInput',
  Pipe = 'Pipe',
  OpenInterpolation = 'OpenInterpolation',
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
export const completeSourceAfter: string[] = [Tokens.Source];

export const comparatorList = [
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

export const parserList = [
  {
    label: 'json_path',
    type: 'operatorKeyword',
    apply: " json_path '",
  },
];

export enum SupportedEditors {
  Selector = 'selector-editor',
  Interpolation = 'interpolation-editor',
  Expression = 'expression-editor',
  CurlCommand = 'curlCommand-editor',
  Definition = 'definition-editor',
}

export const operatorList = [
  {type: 'operatorKeyword', label: '+', apply: '+'},
  {type: 'operatorKeyword', label: '-', apply: '-'},
  {type: 'operatorKeyword', label: '*', apply: '*'},
  {type: 'operatorKeyword', label: '/', apply: '/'},
  {type: 'operatorKeyword', label: '%', apply: '%'},
  {type: 'operatorKeyword', label: '^', apply: '^'},
];

export const SourceByEditorType = {
  [SupportedEditors.Expression]: [
    {
      label: 'env:',
      type: 'variableName',
    },
    {
      label: 'attr:',
      type: 'variableName',
    },
  ],
  [SupportedEditors.Interpolation]: [
    {
      label: 'env:',
      type: 'variableName',
    },
  ],
} as const;
