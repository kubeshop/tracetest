import {PseudoSelector} from 'constants/Operator.constants';
import SelectorService from 'services/Selector.service';
import {TSpanSelector} from 'types/Assertion.types';
import {filterList as filters} from 'utils/Tests';
import {TCompareOperatorSymbol} from '../../types/Operator.types';

describe('AssertionService', () => {
  describe('getSelectorString', () => {
    test('empty selectorList', () => {
      const result = SelectorService.getSelectorString([]);
      expect(result).toBe('');
    });

    test('single selectorList', () => {
      const selectorList = [{operator: '=' as TCompareOperatorSymbol, key: 'service.name', value: 'pokeshop'}];
      const result = SelectorService.getSelectorString(selectorList);
      expect(result).toStrictEqual(`span[service.name = "pokeshop"]`);
    });
    test('single selectorList with pseudoSelector', () => {
      const selectorList = [{operator: '=' as TCompareOperatorSymbol, key: 'service.name', value: 'pokeshop'}];
      const result = SelectorService.getSelectorString(selectorList, {selector: PseudoSelector.LAST});
      expect(result).toStrictEqual(`span[service.name = "pokeshop"]:last`);
    });
    test('single selectorList with nTH pseudoSelector', () => {
      const selectorList = [{operator: '=' as TCompareOperatorSymbol, key: 'service.name', value: 'pokeshop'}];
      const result = SelectorService.getSelectorString(selectorList, {selector: PseudoSelector.NTH, number: 3});
      expect(result).toStrictEqual(`span[service.name = "pokeshop"]:nth_child(3)`);
    });

    test('double selectorList', () => {
      const result = SelectorService.getSelectorString([
        {
          operator: '=',
          key: 'service.name',
          value: 'pokeshop',
        },
        {
          operator: 'contains',
          key: 'tracetest.span.type',
          value: 'http',
        },
      ]);
      expect(result).toStrictEqual(`span[service.name = "pokeshop"  tracetest.span.type contains "http"]`);
    });

    test('with a number value', () => {
      const result = SelectorService.getSelectorString([
        {
          operator: '=',
          key: 'response.body',
          value: 33 as unknown as string,
        },
      ]);
      expect(result).toStrictEqual(`span[response.body = 33]`);
    });

    test('with a JSON value', () => {
      const result = SelectorService.getSelectorString([
        {
          operator: '=',
          key: 'response.body',
          value: JSON.stringify({count: 33}),
        },
      ]);
      expect(result).toStrictEqual(`span[response.body = {\\\"count\\\":33}]`);
    });
  });

  describe('getSpanSelectorList', () => {
    it('should get a list of selector objects from the selector string', () => {
      const selectorString = 'span[service.name = "pokeshop"]';
      const result = SelectorService.getSpanSelectorList(selectorString);
      expect(result).toStrictEqual([
        {
          key: 'service.name',
          operator: '=',
          value: 'pokeshop',
        },
      ]);
    });

    it('should get a list of selector objects from the selector string without quotes', () => {
      const selectorString = 'span[service.name = pokeshop]';
      const result = SelectorService.getSpanSelectorList(selectorString);
      expect(result).toStrictEqual([
        {
          key: 'service.name',
          operator: '=',
          value: 'pokeshop',
        },
      ]);
    });

    it('should get a list of selector objects from the selector string with multiple filters', () => {
      const selectorString = 'span[service.name = "pokeshop"  tracetest.span.type = "http"]';
      const result = SelectorService.getSpanSelectorList(selectorString);
      expect(result).toStrictEqual([
        {
          key: 'service.name',
          operator: '=',
          value: 'pokeshop',
        },
        {
          key: 'tracetest.span.type',
          operator: '=',
          value: 'http',
        },
      ]);
    });
  });

  describe('getSpanSelectorListFromStructure', () => {
    it('should get a list of selector objects from the selector string', () => {
      const firstFilterOnly = [...filters].splice(0, 1);
      const result = SelectorService.getSpanSelectorListFromStructure([{filters: firstFilterOnly}]);
      expect(result).toStrictEqual([
        {
          key: 'service.name',
          operator: '=',
          value: 'pokeshop',
        },
      ]);
    });

    it('should get a list of selector objects from the selector string with multiple filters', () => {
      const result = SelectorService.getSpanSelectorListFromStructure([{filters}]);
      expect(result).toStrictEqual([
        {
          key: 'service.name',
          operator: '=',
          value: 'pokeshop',
        },
        {
          key: 'tracetest.span.type',
          operator: '=',
          value: 'http',
        },
      ]);
    });
  });

  describe('getPseudoSelectorString', () => {
    it('should get a pseudo selector object from a selector string', () => {
      const selector = 'span[service.name="pokeshop" tracetest.span.type="http"]';
      const result = SelectorService.getPseudoSelector(selector);
      expect(result?.selector).toStrictEqual(undefined);
    });

    it('should get a pseudo selector object from a selector string', () => {
      const selector = 'span[service.name="pokeshop" tracetest.span.type="http"]:first';
      const result = SelectorService.getPseudoSelector(selector);
      expect(result?.selector).toStrictEqual(PseudoSelector.FIRST);
    });

    it('should get a selector object from a nth_child selector', () => {
      const selector = 'span[service.name="pokeshop" tracetest.span.type="http"]:nth_child(2)';
      const result = SelectorService.getPseudoSelector(selector);
      expect(result).toStrictEqual({
        selector: PseudoSelector.NTH,
        number: 2,
      });
    });
  });
  describe('getPseudoSelectorFromStructure', () => {
    it('should get a pseudo selector object from a selector string', () => {
      const result = SelectorService.getPseudoSelectorFromStructure([{filters}]);
      expect(result?.selector).toStrictEqual(undefined);
    });

    it('should get a pseudo selector object from a selector string', () => {
      const pseudoClass = {name: 'first'};
      const result = SelectorService.getPseudoSelectorFromStructure([{filters, pseudoClass}]);
      expect(result?.selector).toStrictEqual(PseudoSelector.FIRST);
    });

    it('should get a selector object from a nth_child selector', () => {
      const pseudoClass = {name: 'nth_child', argument: 2};
      const result = SelectorService.getPseudoSelectorFromStructure([{filters, pseudoClass}]);
      expect(result).toStrictEqual({
        selector: PseudoSelector.NTH,
        number: 2,
      });
    });
  });
  describe('validateSelector', () => {
    it('should return true for a valid selector', async () => {
      const selector: TSpanSelector = {
        key: 'service.name',
        operator: '=',
        value: 'pokeshop',
      };

      const result = await SelectorService.validateSelector([], false, [], [selector]);
      expect(result).toEqual(true);
    });

    it('should return true when editing and the selector match with the initial', async () => {
      const selector: TSpanSelector = {
        key: 'service.name',
        operator: '=',
        value: 'pokeshop',
      };

      const selectorString = 'span[service.name = "pokeshop"]';

      const result = await SelectorService.validateSelector([selectorString], true, [selector], [selector]);
      expect(result).toEqual(true);
    });

    it('should return an error when editing and the selector does not match with the initial', done => {
      expect.assertions(2);
      const selector: TSpanSelector = {
        key: 'service.name',
        operator: '=',
        value: 'pokeshop',
      };

      const selectorString = 'span[service.name = "pokeshop"]';

      SelectorService.validateSelector([selectorString], true, [], [selector]).catch((error: Error) => {
        expect(error).toBeInstanceOf(Error);
        expect(error.message).toBe('Selector already exists');
        done();
      });
    });
  });

  it('should return an error when not editing and the selector is already part of the list', done => {
    expect.assertions(2);
    const selector: TSpanSelector = {
      key: 'service.name',
      operator: '=',
      value: 'pokeshop',
    };

    const selectorString = 'span[service.name = "pokeshop"]';

    SelectorService.validateSelector([selectorString], false, [selector], [selector]).catch((error: Error) => {
      expect(error).toBeInstanceOf(Error);
      expect(error.message).toBe('Selector already exists');

      done();
    });
  });
});
