import {RootState} from '../../redux/store';
import {ITestDefinitionState} from '../../types/TestSpecs.types';
import TestDefinitionSelectors from '../TestDefinition.selectors';

describe('TestDefinitionSelectors', () => {
  const initialTestDefinitionState = {
    initialDefinitionList: [],
    definitionList: [],
    changeList: [],
    isLoading: false,
    isInitialized: false,
    selectedAssertion: undefined,
    isDraftMode: false,
  };

  describe('selectDefinitionList', () => {
    it('should return empty', () => {
      const result = TestDefinitionSelectors.selectDefinitionList({
        testDefinition: {...initialTestDefinitionState} as ITestDefinitionState,
      } as RootState);
      expect(result).toStrictEqual([]);
    });
  });
  describe('selectDefinitionSelectorList', () => {
    it('should return array of definitionList selectors', () => {
      const result = TestDefinitionSelectors.selectDefinitionSelectorList({
        testDefinition: {definitionList: [{selector: 'hello'}]} as ITestDefinitionState,
      } as RootState);
      expect(result).toStrictEqual(['hello']);
    });
  });
  describe('selectIsSelectorExist', () => {
    it('should return true', () => {
      const result = TestDefinitionSelectors.selectIsSelectorExist(
        {
          testDefinition: {definitionList: [{selector: 'hello'}]} as ITestDefinitionState,
        } as RootState,
        'hello'
      );
      expect(result).toStrictEqual(true);
    });
  });
  describe('selectIsLoading', () => {
    it('should return false', () => {
      const result = TestDefinitionSelectors.selectIsLoading({
        testDefinition: {isLoading: false} as ITestDefinitionState,
      } as RootState);
      expect(result).toStrictEqual(false);
    });
  });
  describe('selectIsInitialized', () => {
    it('should return false', () => {
      const result = TestDefinitionSelectors.selectIsInitialized({
        testDefinition: {isInitialized: false} as ITestDefinitionState,
      } as RootState);
      expect(result).toStrictEqual(false);
    });
  });
  describe('selectAssertionResults', () => {
    it('should return assertionResults', () => {
      const assertionResults = {allPassed: false, resultList: [], results: undefined};
      const result = TestDefinitionSelectors.selectAssertionResults({
        testDefinition: {
          ...initialTestDefinitionState,
          assertionResults,
        } as ITestDefinitionState,
      } as RootState);
      expect(result).toStrictEqual(assertionResults);
    });
  });
  describe('selectDefinitionBySelector', () => {
    it('should return assertionResults', () => {
      const selector = 'selector';
      const definition = {selector};
      const result = TestDefinitionSelectors.selectDefinitionBySelector(
        {
          testDefinition: {definitionList: [definition]} as ITestDefinitionState,
        } as RootState,
        selector
      );
      expect(result).toStrictEqual(definition);
    });
  });
  describe('selectSelectedAssertion', () => {
    it('should return false', () => {
      const selectedAssertion = 'thisAssertion';
      const result = TestDefinitionSelectors.selectSelectedAssertion({
        testDefinition: {selectedAssertion} as ITestDefinitionState,
      } as RootState);
      expect(result).toStrictEqual(selectedAssertion);
    });
  });
  describe('selectAssertionResultsBySpan', () => {
    it('should return empty object', () => {
      const assertionResults = {allPassed: false, resultList: [], results: undefined};
      const result = TestDefinitionSelectors.selectAssertionResultsBySpan(
        {
          testDefinition: {
            ...initialTestDefinitionState,
            assertionResults,
          } as ITestDefinitionState,
        } as RootState,
        'spanId'
      );
      expect(result).toStrictEqual({});
    });
  });
  describe('selectAssertionResultsBySpan', () => {
    it('should return false', () => {
      const result = TestDefinitionSelectors.selectIsDraftMode({
        testDefinition: {isDraftMode: false} as ITestDefinitionState,
      } as RootState);
      expect(result).toStrictEqual(false);
    });
  });
});
