import {isNumber} from 'lodash';
import {CompareOperator} from '../constants/Operator.constants';
import {TSpanSelector} from '../types/Assertion.types';
import {escapeString, isJson} from '../utils/Common';
import OperatorService from './Operator.service';

const getValue = (value: string): string => {
  if (isNumber(value)) {
    return value;
  }

  if (isJson(value)) {
    return escapeString(value);
  }

  return value;
};

const getFilters = (selectors: TSpanSelector[]) =>
  selectors.map(
    ({key, operator = CompareOperator.EQUALS, value}) =>
      `${key}${OperatorService.getOperatorSymbol(operator)}${getValue(value)}`
  );

const AssertionService = () => ({
  getSelectorString(selectorList: TSpanSelector[]): string {
    return selectorList.length ? `span[${getFilters(selectorList).join(' ')}]` : '';
  },

  getSpanSelectorList(selectorString: string): TSpanSelector[] {
    console.log('@@selectorString', selectorString);
    return [];
  },
});

export default AssertionService();
