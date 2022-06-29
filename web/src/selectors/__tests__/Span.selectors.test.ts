import SpanSelectors from 'selectors/Span.selectors';
import SpanMock from 'models/__mocks__/Span.mock';
import {RootState} from 'redux/store';
import {TAssertionResultEntry} from 'types/Assertion.types';
import {ISpanState, TSpan} from 'types/Span.types';

describe('SpanSelectors', () => {
  describe('selectAffectedSpans', () => {
    it('should return affectedSpans when selectedAssertion not present', () => {
      const affectedSpans = ['pokeshop'];
      const result = SpanSelectors.selectAffectedSpans({
        spans: {affectedSpans},
        testDefinition: {},
      } as RootState);
      expect(result).toBe(affectedSpans);
    });
    it('should return affectedSpans when selector matches', () => {
      const affectedSpans = ['pokeshop'];
      const selector = `span[tracetest.span.type="http"]`;
      const result = SpanSelectors.selectAffectedSpans({
        spans: {affectedSpans},
        testDefinition: {
          assertionResults: {resultList: [{selector} as TAssertionResultEntry]},
          selectedAssertion: selector,
        },
      } as RootState);
      expect(result).toBe(affectedSpans);
    });
    it('should return affectedSpans when selector does not matches', () => {
      const selector = `span[tracetest.span.type="http"]`;
      const selectedAssertion = `span[tracetest.span.type="gRPC"]`;
      const result = SpanSelectors.selectAffectedSpans({
        spans: {affectedSpans: ['pokeshop']},
        testDefinition: {
          assertionResults: {resultList: [{selector} as TAssertionResultEntry]},
          selectedAssertion,
        },
      } as RootState);
      expect(result).toStrictEqual([]);
    });
  });
  describe('selectSelectedSpan', () => {
    it('should return span', () => {
      const selectedSpan: TSpan = SpanMock.model();
      const result = SpanSelectors.selectSelectedSpan({
        spans: {selectedSpan} as ISpanState,
      } as RootState);
      expect(result).toBe(selectedSpan);
    });
  });

  describe('selectFocusedSpan', () => {
    it('should return focusedSpan', () => {
      const focusedSpan = 'string';
      const result = SpanSelectors.selectFocusedSpan({
        spans: {focusedSpan} as ISpanState,
      } as RootState);
      expect(result).toBe(focusedSpan);
    });
  });
  describe('selectMatchedSpans', () => {
    it('should return matchedSpans', () => {
      const matchedSpans = ['string'];
      const result = SpanSelectors.selectMatchedSpans({
        spans: {matchedSpans} as ISpanState,
      } as RootState);
      expect(result).toBe(matchedSpans);
    });
  });

  describe('selectSearchText', () => {
    it('should return searchText', () => {
      const searchText = 'string';
      const result = SpanSelectors.selectSearchText({
        spans: {searchText} as ISpanState,
      } as RootState);
      expect(result).toBe(searchText);
    });
  });
});
