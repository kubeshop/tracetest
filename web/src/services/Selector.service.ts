import {isNumber} from 'lodash';
import {PseudoSelector} from '../constants/Operator.constants';
import {TPseudoSelector, TSpanSelector} from '../types/Assertion.types';
import {TFilter, TStructure} from '../types/Common.types';
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

const nthChildNumberRegex = /\((.*)\)/i;
const selectorRegex = /span\[(.*)\]/i;
const operationRegex = /\s?([=<]+|contains)\s?/;

const getFilters = (selectors: TSpanSelector[]) =>
  selectors.map(({key, operator, value}) => `${key} ${operator} ${getValue(value)}`);

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

function flattenStructureFilters(structure: TStructure): TFilter[] {
  return [
    ...(structure?.filters || []),
    ...(structure?.childSelector ? flattenStructureFilters(structure.childSelector) : []),
  ];
}

const SelectorService = () => ({
  getSelectorString(selectorList: TSpanSelector[], pseudoSelector?: TPseudoSelector): string {
    return selectorList.length
      ? `span[${getFilters(selectorList).join(' ')}]${getPseudoSelectorString(pseudoSelector)}`
      : '';
  },

  getSpanSelectorList(selectorString: string): TSpanSelector[] {
    const [matchString] = (selectorString.match(selectorRegex) || []).reverse();

    if (!matchString) return [];

    return matchString.split('  ').reduce<TSpanSelector[]>((list, operation) => {
      if (!operation) return list;
      const [key, operator, value] = operation.split(operationRegex);

      const spanSelector = {
        key: key.trim(),
        operator: operator.trim() as TCompareOperatorSymbol,
        value: getCleanValue(value.trim()),
      };

      return list.concat([spanSelector]);
    }, []);
  },
  getSpanSelectorListFromStructure(structures: TStructure[]): TSpanSelector[] {
    return structures
      .flatMap(a => flattenStructureFilters(a))
      .map(r => ({
        key: r?.property || '',
        operator: r?.operator as TCompareOperatorSymbol,
        value: r?.value || '',
      }));
  },

  getPseudoSelectorFromStructure(structure: TStructure[]): TPseudoSelector | undefined {
    const pseudoSelector = structure[0]?.pseudoClass || undefined;
    return pseudoSelector
      ? {selector: `:${pseudoSelector.name}` as PseudoSelector, number: pseudoSelector.argument}
      : undefined;
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

  validateSelector(
    definitionSelectorList: string[],
    isEditing: boolean,
    initialSelectorList: TSpanSelector[],
    selectorList: TSpanSelector[],
    initialPseudoSelector?: TPseudoSelector,
    pseudoSelector?: TPseudoSelector
  ): Promise<boolean> {
    const initialSelectorString = this.getSelectorString(initialSelectorList, initialPseudoSelector);
    const selectorString = this.getSelectorString(selectorList, pseudoSelector);

    if (!definitionSelectorList.includes(selectorString) || (isEditing && initialSelectorString === selectorString))
      return Promise.resolve(true);
    return Promise.reject(new Error('Selector already exists'));
  },

  getIsAdvancedSelector(selector: string): boolean {
    const matches = (selector.match(/span\[/g) || []).length;
    const hasOrOperator = selector.includes('],');

    return matches > 1 || hasOrOperator;
  },
});

export default SelectorService();
