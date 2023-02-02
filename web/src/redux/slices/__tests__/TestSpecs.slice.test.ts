import faker from '@faker-js/faker';
import TestSpecs, {TTestSpecEntry} from 'models/TestSpecs.model';
import AssertionResultsMock from '../../../models/__mocks__/AssertionResults.mock';
import TestDefinitionMock from '../../../models/__mocks__/TestSpecs.mock';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';
import Reducer, {
  addSpec,
  assertionResultsToSpecs,
  initSpecs,
  initialState,
  removeSpec,
  reset,
  resetSpecs,
  revertSpec,
  setAssertionResults,
  updateSpec,
  setSelectedSpec,
} from '../TestSpecs.slice';

const {specs} = TestDefinitionMock.model();

const specSelector = `span[http.status_code] = "304"]`;

const spec: TTestSpecEntry = {
  selector: specSelector,
  isDraft: true,
  assertions: new Array(faker.datatype.number({min: 2, max: 10}))
    .fill(null)
    .map(() => `${faker.datatype.string(10)} = "${faker.datatype.string(10)}"`),
  originalSelector: specSelector,
  name: '',
};

const state = {...initialState, specs};

describe('TestSpecs slice', () => {
  it('should return the initial state', () => {
    expect(Reducer(undefined, {type: 'any-action'})).toEqual(initialState);
  });

  describe('initSpecs', () => {
    it('should handle triggering the action', () => {
      const assertionResults = AssertionResultsMock.model();
      expect(Reducer(undefined, initSpecs({assertionResults, specs: TestSpecs({})}))).toEqual({
        ...initialState,
        initialSpecs: assertionResultsToSpecs(assertionResults, TestSpecs({})),
        specs: assertionResultsToSpecs(assertionResults, TestSpecs({})),
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

  describe('resetSpecs', () => {
    it('should handle resetting the specs', () => {
      const result = Reducer(state, resetSpecs());

      expect(result).toEqual({
        ...state,
        specs: [],
      });
    });
  });

  describe('reset', () => {
    it('should handle resetting the whole state', () => {
      const result = Reducer(state, reset());

      expect(result).toEqual(initialState);
    });
  });

  describe('spec CRUD', () => {
    it('should handle the add spec action', () => {
      expect(Reducer(state, addSpec({spec}))).toEqual({
        ...initialState,
        specs: [...specs, spec],
        isDraftMode: true,
      });
    });

    it('should handle the update spec action', () => {
      const result = Reducer(
        state,
        updateSpec({
          selector: specs[0].selector,
          spec,
        })
      );

      expect(result.specs[0]).toEqual({...spec, originalSelector: undefined});
      expect(result).toEqual({
        ...initialState,
        specs: [{...spec, originalSelector: undefined}, ...specs.slice(1, specs.length)],
        isDraftMode: true,
      });
    });

    it('should handle the revert spec action', () => {
      const initialSelector = 'span[http.status_code = "204"]';
      const list = [
        {
          ...spec,
          originalSelector: initialSelector,
        },
      ];
      const result = Reducer(
        {
          ...state,
          initialSpecs: list,
          specs: list,
        },
        revertSpec({
          originalSelector: initialSelector,
        })
      );

      expect(result.specs[0].originalSelector).toEqual(initialSelector);
      expect(result.specs[0].selector).toEqual(specSelector);
    });

    it('should handle the remove spec action', () => {
      const result = Reducer(
        state,
        removeSpec({
          selector: specs[0].selector,
        })
      );

      expect(result).toEqual({
        ...initialState,
        specs: [{...specs[0], isDraft: true, isDeleted: true}, ...specs.slice(1, specs.length)],
        isDraftMode: true,
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
      expect(result.initialSpecs).toEqual(assertionResultsToSpecs(run.result, TestSpecs({})));
    });
  });

  it('should handle on pending publish', () => {
    const result = Reducer(initialState, {
      type: 'testDefinition/publish/pending',
    });

    expect(result.isDraftMode).toBeFalsy();
  });

  describe('setSelectedSpec', () => {
    it('should handle on setSelectedAssertion', () => {
      const assertionResultEntry = {
        id: faker.datatype.uuid(),
        selector: faker.random.word(),
        originalSelector: faker.random.word(),
        spanIds: ['12345', '67890'],
        resultList: [],
      };

      const result = Reducer(initialState, setSelectedSpec(assertionResultEntry));

      expect(result.selectedSpec).toEqual(assertionResultEntry.selector);
    });

    it('should handle on setSelectedSpec with empty value', () => {
      const result = Reducer({...initialState, selectedSpec: '12345'}, setSelectedSpec());

      expect(result.selectedSpec).toBeUndefined();
    });
  });
});
