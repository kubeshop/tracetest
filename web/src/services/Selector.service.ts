import {isNumber} from 'lodash';
import {PseudoSelector} from '../constants/Operator.constants';
import {TPseudoSelector, TSpanSelector} from '../types/Assertion.types';
import {TCompareOperatorSymbol} from '../types/Operator.types';
import {escapeString, isJson} from '../utils/Common';

const getValue = (value: string): string => {
  if (isNumber(value)) {
    return value;
  }

  if (isJson(value)) {
    return escapeString(value);
  }

  return `"${value}"`;
};

const selectorRegex = /span\[(.*)\]/i;
const nthChildNumberRegex = /\((.*)\)/i;
const operationRegex = /([=]+|contains)/;

const getFilters = (selectors: TSpanSelector[]) =>
  selectors.map(({key, operator, value}) => `${key}${operator}${getValue(value)}`);

const getCleanValue = (value: string): string => {
  if (value.includes('"')) {
    return value.replace(/"/g, '');
  }

  return value;
};

const getPseudoSelectorString = (pseudoSelector?: TPseudoSelector): string => {
  if (!pseudoSelector) return '';

  const {selector, number} = pseudoSelector;

  if (selector === PseudoSelector.NTH) {
    return `${pseudoSelector.selector}(${number})`;
  }

  return selector;
};

const SelectorService = () => ({
  getSelectorString(selectorList: TSpanSelector[], pseudoSelector?: TPseudoSelector): string {
    return selectorList.length
      ? `span[${getFilters(selectorList).join(' ')}]${getPseudoSelectorString(pseudoSelector)}`
      : '';
  },

  getSpanSelectorList(selectorString: string): TSpanSelector[] {
    const [matchString] = (selectorString.match(selectorRegex) || []).reverse();

    if (!matchString) return [];

    const selectorList = matchString.split(' ').map(operation => {
      const [key, operator, value] = operation.split(operationRegex);

      return {key, operator: operator as TCompareOperatorSymbol, value: getCleanValue(value)};
    });

    return selectorList;
  },

  getPseudoSelector(selectorString: string): TPseudoSelector | undefined {
    const index = selectorString.indexOf(']:');
    if (index === -1) return;

    const pseudoSelector = selectorString.substring(selectorString.indexOf(']:') + 1);

    if (pseudoSelector.includes('nth_child')) {
      const [number] = (pseudoSelector.match(nthChildNumberRegex) || []).reverse();

      return {selector: PseudoSelector.NTH, number: Number(number)};
    }

    return pseudoSelector
      ? {
          selector: pseudoSelector as PseudoSelector,
        }
      : undefined;
  },
});

export default SelectorService();
