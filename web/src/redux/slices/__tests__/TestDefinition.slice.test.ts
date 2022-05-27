import faker from '@faker-js/faker';
import AssertionMock from '../../../models/__mocks__/Assertion.mock';
import AssertionResultsMock from '../../../models/__mocks__/AssertionResults.mock';
import TestDefinitionMock from '../../../models/__mocks__/TestDefinition.mock';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';
import {TTestDefinitionEntry} from '../../../types/TestDefinition.types';
import Reducer, {
  addDefinition,
  initDefinitionList,
  resetDefinitionList,
  removeDefinition,
  reset,
  setAssertionResults,
  initialState,
  assertionResultsToDefinitionList,
  updateDefinition,
  revertDefinition,
} from '../TestDefinition.slice';

const {definitionList} = TestDefinitionMock.model();
const definition: TTestDefinitionEntry = {
  selector: `span[http.status_code] = "304"]`,
  isDraft: true,
  assertionList: new Array(faker.datatype.number({min: 2, max: 10})).fill(null).map(() => AssertionMock.model()),
};

const state = {
  ...initialState,
  definitionList,
};

describe('TestDefinitionReducer', () => {
  it('should return the initial state', () => {
    expect(Reducer(undefined, {type: 'any-action'})).toEqual(initialState);
  });

  describe('initDefinitionList', () => {
    it('should handle triggering the action', () => {
      const assertionResults = AssertionResultsMock.model();
      expect(Reducer(undefined, initDefinitionList({assertionResults}))).toEqual({
        ...initialState,
        initialDefinitionList: assertionResultsToDefinitionList(assertionResults),
        definitionList: assertionResultsToDefinitionList(assertionResults),
        isInitialized: true,
      });
    });
  });

  describe('setAssertionResults', () => {
    it('should handle triggering the action', () => {
      const assertionResults = AssertionResultsMock.model();
      const result = Reducer(state, setAssertionResults(assertionResults));

      expect(result).toEqual({
        ...state,
        assertionResults,
      });
    });
  });

  describe('resetDefinitionList', () => {
    it('should handle resetting the definition list', () => {
      const result = Reducer(state, resetDefinitionList());

      expect(result).toEqual({
        ...state,
        definitionList: [],
      });
    });
  });

  describe('reset', () => {
    it('should handle resetting the whole state', () => {
      const result = Reducer(state, reset());

      expect(result).toEqual(initialState);
    });
  });

  describe('definition cRUD', () => {
    it('should handle the add definition action', () => {
      expect(Reducer(state, addDefinition({definition}))).toEqual({
        ...initialState,
        definitionList: [...definitionList, definition],
      });
    });

    it('should handle the updating definition action', () => {
      const result = Reducer(
        state,
        updateDefinition({
          selector: definitionList[0].selector,
          definition,
        })
      );

      expect(result.definitionList[0]).toEqual(definition);
      expect(result).toEqual({
        ...initialState,
        definitionList: [definition, ...definitionList.slice(1, definitionList.length)],
      });
    });

    it('should handle the revert definition action', () => {
      const initialSelector = 'span[http.status_code] = "204"]';
      const result = Reducer(
        {
          ...state,
          initialDefinitionList: [
            {
              ...state.definitionList[0],
              isDraft: true,
              selector: initialSelector,
            },
          ],
        },
        revertDefinition({
          index: 0,
        })
      );

      expect(result.definitionList[0].selector).toEqual(initialSelector);
    });

    it('should handle the remove definition action', () => {
      const result = Reducer(
        state,
        removeDefinition({
          selector: definitionList[0].selector,
        })
      );

      expect(result).toEqual({
        ...initialState,
        definitionList: [
          {...definitionList[0], isDraft: true, isDeleted: true},
          ...definitionList.slice(1, definitionList.length),
        ],
      });
    });
  });

  describe('loading', () => {
    it('should handle updating the loading state to false', () => {
      const result = Reducer(initialState, {type: 'testDefinition/dryRun/fulfilled'});

      expect(result.isLoading).toEqual(false);
    });

    it('should handle updating the loading state to true', () => {
      const result = Reducer(initialState, {type: 'testDefinition/dryRun/pending'});

      expect(result.isLoading).toEqual(true);
    });
  });

  describe('dryRun', () => {
    it('should handle on fulfilled dry run', () => {
      const run = TestRunMock.model();
      const result = Reducer(initialState, {
        type: 'testDefinition/dryRun/fulfilled',
        payload: run.result,
      });

      expect(result.assertionResults).toEqual(run.result);
    });
  });

  describe('publish', () => {
    it('should handle on fulfilled publish', () => {
      const run = TestRunMock.model();
      const result = Reducer(initialState, {
        type: 'testDefinition/publish/fulfilled',
        payload: run,
      });

      expect(result.assertionResults).toEqual(run.result);
      expect(result.initialDefinitionList).toEqual(assertionResultsToDefinitionList(run.result));
    });
  });
});
