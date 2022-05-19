import {IItemSelector} from '../../types/Assertion.types';

export function newSelectorLogic(selectorList: IItemSelector[]): string {
  let result = '';
  for (let i = 0; i < selectorList.length; i += 1) {
    if (i === 0) result += `span[`;
    const selector = selectorList[i];
    result += selector.propertyName;
    result += selector.operator ? ` ${selector.operator.toLowerCase()} ` : '=';
    // add quotes if value is a string
    if (selector.valueType === 'stringValue') result += `'`;
    result += selector.value;
    // add quotes if value is a string
    if (selector.valueType === 'stringValue') result += `'`;
    if (i === selectorList.length - 1) result += `]`;
    // adds space between selector if not the last selector
    if (i !== selectorList.length - 1) result += ` `;
  }
  return result;
}
