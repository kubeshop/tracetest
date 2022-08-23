import faker from '@faker-js/faker';
import SpanMock from 'models/__mocks__/Span.mock';
import Reducer, {
  initialState,
  setAffectedSpans,
  setFocusedSpan,
  setSelectedSpan,
  clearAffectedSpans,
  setMatchedSpans,
  setSearchText,
} from '../Span.slice';
import {setSelectedAssertion} from '../TestDefinition.slice';

describe('Span.slice', () => {
  it('should return the initial state', () => {
    expect(Reducer(undefined, {type: 'any-action'})).toEqual(initialState);
  });

  describe('setAffectedSpans', () => {
    it('should handle triggering the action', () => {
      expect(Reducer(undefined, setAffectedSpans({spanIds: ['12345', '567890']}))).toEqual({
        ...initialState,
        affectedSpans: ['12345', '567890'],
        focusedSpan: '12345',
      });
    });

    it('should handle an empty array', () => {
      expect(Reducer(undefined, setAffectedSpans({spanIds: []}))).toEqual(initialState);
    });
  });

  describe('setFocusedSpan', () => {
    it('should handle triggering the action', () => {
      expect(Reducer(undefined, setFocusedSpan({spanId: '12345'}))).toEqual({
        ...initialState,
        focusedSpan: '12345',
      });
    });
  });

  describe('setSelectedSpan', () => {
    it('should handle triggering the action', () => {
      const span = SpanMock.model();
      expect(Reducer(undefined, setSelectedSpan({span}))).toEqual({
        ...initialState,
        selectedSpan: span,
      });
    });
  });

  describe('clearAffectedSpans', () => {
    it('should handle triggering the action', () => {
      expect(
        Reducer({...initialState, focusedSpan: '', affectedSpans: ['12345', '67890']}, clearAffectedSpans())
      ).toEqual(initialState);
    });
  });

  describe('setSelectedAssertion side effect', () => {
    it('should handle triggering the action', () => {
      const assertionResultEntry = {
        id: faker.datatype.uuid(),
        selector: faker.random.word(),
        originalSelector: faker.random.word(),
        spanIds: ['12345', '67890'],
        resultList: [],
      };
      expect(
        Reducer(
          {...initialState, focusedSpan: '', affectedSpans: ['12345', '67890']},
          setSelectedAssertion(assertionResultEntry)
        )
      ).toEqual({
        ...initialState,
        affectedSpans: ['12345', '67890'],
        focusedSpan: '12345',
      });
    });

    it('should handle removing the selected assertion', () => {
      expect(
        Reducer(
          {
            ...initialState,
            affectedSpans: ['12345', '67890'],
            focusedSpan: '12345',
          },
          setSelectedAssertion()
        )
      ).toEqual(initialState);
    });
  });

  describe('setMatchedSpans', () => {
    it('should handle triggering the action', () => {
      expect(Reducer(undefined, setMatchedSpans({spanIds: ['12345', '567890']}))).toEqual({
        ...initialState,
        matchedSpans: ['12345', '567890'],
      });
    });
  });

  describe('setSearchText', () => {
    it('should handle triggering the action', () => {
      expect(Reducer(undefined, setSearchText({searchText: 'http'}))).toEqual({
        ...initialState,
        searchText: 'http',
      });
    });
  });
});
