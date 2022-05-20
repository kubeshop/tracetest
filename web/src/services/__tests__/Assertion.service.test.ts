import {CompareOperator} from '../../constants/Operator.constants';
import AssertionService from '../Assertion.service';

describe('AssertionService', () => {
  describe('getSelectorString', () => {
    test('empty selectorList', () => {
      const result = AssertionService.getSelectorString([]);
      expect(result).toBe('');
    });

    test('single selectorList', () => {
      const result = AssertionService.getSelectorString([
        {
          operator: CompareOperator.EQUALS,
          key: 'service.name',
          value: 'pokeshop',
        },
      ]);
      expect(result).toStrictEqual(`span[service.name="pokeshop"]`);
    });

    test('double selectorList', () => {
      const result = AssertionService.getSelectorString([
        {
          operator: CompareOperator.EQUALS,
          key: 'service.name',
          value: 'pokeshop',
        },
        {
          operator: CompareOperator.CONTAINS,
          key: 'tracetest.span.type',
          value: 'http',
        },
      ]);
      expect(result).toStrictEqual(`span[service.name="pokeshop" tracetest.span.type contains "http"]`);
    });
  });
});
