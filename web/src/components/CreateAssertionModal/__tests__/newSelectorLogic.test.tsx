import {CompareOperator} from '../../../constants/Operator.constants';
import {LOCATION_NAME} from '../../../constants/Span.constants';
import {newSelectorLogic} from '../newSelectorLogic';

describe('newSelectorLogic', () => {
  test('empty selectorList', () => {
    const result = newSelectorLogic([]);
    expect(result).toBe('');
  });
  test('single selectorList', () => {
    const result = newSelectorLogic([
      {
        locationName: LOCATION_NAME.SPAN,
        operator: undefined,
        propertyName: 'service.name',
        value: 'pokeshop',
        valueType: 'stringValue',
      },
    ]);
    expect(result).toStrictEqual(`span[service.name='pokeshop']`);
  });
  test('double selectorList', () => {
    const result = newSelectorLogic([
      {
        locationName: LOCATION_NAME.SPAN,
        operator: undefined,
        propertyName: 'service.name',
        value: 'pokeshop',
        valueType: 'stringValue',
      },
      {
        locationName: LOCATION_NAME.SPAN,
        operator: CompareOperator.CONTAINS,
        propertyName: 'tracetest.span.type',
        value: 'http',
        valueType: 'stringValue',
      },
    ]);
    expect(result).toStrictEqual(`span[service.name='pokeshop' tracetest.span.type contains 'http']`);
  });
});
