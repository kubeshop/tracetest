import {RootState} from '../../redux/store';
import {ITestSpecsState} from '../../types/TestSpecs.types';
import TestSpecsSelectors from '../TestSpecs.selectors';

describe('TestSpecsSelectors', () => {
  const initialTestSpecsState = {
    initialSpecs: [],
    specs: [],
    changeList: [],
    isLoading: false,
    isInitialized: false,
    selectedSpec: undefined,
    isDraftMode: false,
    selectorSuggestions: [],
    prevSelector: '',
  };

  describe('selectSpecs', () => {
    it('should return empty', () => {
      const result = TestSpecsSelectors.selectSpecs({
        testSpecs: {...initialTestSpecsState} as ITestSpecsState,
      } as RootState);
      expect(result).toStrictEqual([]);
    });
  });

  describe('selectSpecSelectorList', () => {
    it('should return array of definitionList selectors', () => {
      const result = TestSpecsSelectors.selectSpecsSelectorList({
        testSpecs: {specs: [{selector: 'hello'}]} as ITestSpecsState,
      } as RootState);
      expect(result).toStrictEqual(['hello']);
    });
  });

  describe('selectIsSelectorExist', () => {
    it('should return true', () => {
      const result = TestSpecsSelectors.selectIsSelectorExist(
        {
          testSpecs: {specs: [{selector: 'hello'}]} as ITestSpecsState,
        } as RootState,
        'hello'
      );
      expect(result).toStrictEqual(true);
    });
  });

  describe('selectIsLoading', () => {
    it('should return false', () => {
      const result = TestSpecsSelectors.selectIsLoading({
        testSpecs: {isLoading: false} as ITestSpecsState,
      } as RootState);
      expect(result).toStrictEqual(false);
    });
  });

  describe('selectIsInitialized', () => {
    it('should return false', () => {
      const result = TestSpecsSelectors.selectIsInitialized({
        testSpecs: {isInitialized: false} as ITestSpecsState,
      } as RootState);
      expect(result).toStrictEqual(false);
    });
  });

  describe('selectAssertionResults', () => {
    it('should return assertionResults', () => {
      const assertionResults = {allPassed: false, resultList: [], results: undefined};
      const result = TestSpecsSelectors.selectAssertionResults({
        testSpecs: {
          ...initialTestSpecsState,
          assertionResults,
        } as ITestSpecsState,
      } as RootState);
      expect(result).toStrictEqual(assertionResults);
    });
  });

  describe('selectSpecBySelector', () => {
    it('should return assertionResults', () => {
      const selector = 'selector';
      const definition = {selector};
      const result = TestSpecsSelectors.selectSpecBySelector(
        {
          testSpecs: {specs: [definition]} as ITestSpecsState,
        } as RootState,
        selector
      );
      expect(result).toStrictEqual(definition);
    });
  });

  describe('selectSelectedSpec', () => {
    it('should return false', () => {
      const selectedSpec = 'thisAssertion';
      const result = TestSpecsSelectors.selectSelectedSpec({
        testSpecs: {selectedSpec} as ITestSpecsState,
      } as RootState);
      expect(result).toStrictEqual(selectedSpec);
    });
  });

  describe('selectIsDraftMode', () => {
    it('should return false', () => {
      const result = TestSpecsSelectors.selectIsDraftMode({
        testSpecs: {isDraftMode: false} as ITestSpecsState,
      } as RootState);
      expect(result).toStrictEqual(false);
    });
  });
});
