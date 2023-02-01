import SpanSelectors from 'selectors/Span.selectors';
import SpanMock from 'models/__mocks__/Span.mock';
import {RootState} from 'redux/store';
import {ISpanState} from 'types/Span.types';
import {TAssertionResultEntry} from '../../models/AssertionResults.model';

describe('SpanSelectors', () => {
  describe('selectMatchedSpans', () => {
    it('should return matchedSpans when selectedAssertion not present', () => {
      const matchedSpans = ['pokeshop'];
      const result = SpanSelectors.selectMatchedSpans({
        spans: {matchedSpans},
        testSpecs: {},
      } as RootState);
      expect(result).toBe(matchedSpans);
    });

    it('should return matchedSpans when selector matches', () => {
      const matchedSpans = ['pokeshop'];
      const selector = `span[tracetest.span.type="http"]`;
      const result = SpanSelectors.selectMatchedSpans({
        spans: {matchedSpans},
        testSpecs: {
          assertionResults: {resultList: [{selector} as TAssertionResultEntry]},
          selectedSpec: selector,
        },
      } as RootState);
      expect(result).toBe(matchedSpans);
    });

    it('should return matchedSpans when selector does not matches', () => {
      const selector = `span[tracetest.span.type="http"]`;
      const selectedSpec = `span[tracetest.span.type="gRPC"]`;
      const result = SpanSelectors.selectMatchedSpans({
        spans: {matchedSpans: ['pokeshop']},
        testSpecs: {
          assertionResults: {resultList: [{selector} as TAssertionResultEntry]},
          selectedSpec,
        },
      } as RootState);
      expect(result).toStrictEqual([]);
    });
  });

  describe('selectSelectedSpan', () => {
    it('should return span', () => {
      const selectedSpan = SpanMock.model();
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
});
