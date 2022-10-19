import {createSlice} from '@reduxjs/toolkit';
import {TAssertionResults} from 'types/Assertion.types';
import {TTestSpecEntry, TTestSpecsSliceActions, ITestSpecsState} from 'types/TestSpecs.types';
import TestSpecsActions from '../actions/TestSpecs.actions';

export const initialState: ITestSpecsState = {
  initialSpecs: [],
  specs: [],
  changeList: [],
  isLoading: false,
  isInitialized: false,
  selectedSpec: undefined,
  isDraftMode: false,
  selectorSuggestions: [],
};

export const assertionResultsToSpecs = (assertionResults: TAssertionResults): TTestSpecEntry[] => {
  return assertionResults.resultList.map(({selector, resultList}) => ({
    isDraft: false,
    isDeleted: false,
    selector,
    originalSelector: selector,
    assertions: resultList.flatMap(({assertion}) => [assertion]),
  }));
};

const testSpecsSlice = createSlice<ITestSpecsState, TTestSpecsSliceActions, 'testSpecs'>({
  name: 'testSpecs',
  initialState,
  reducers: {
    reset() {
      return initialState;
    },
    initSpecs(state, {payload: {assertionResults}}) {
      const specs = assertionResultsToSpecs(assertionResults);

      state.initialSpecs = specs;
      state.specs = specs;
      state.isInitialized = true;
      state.isDraftMode = false;
    },
    setIsInitialized(state, {payload: {isInitialized}}) {
      state.isInitialized = isInitialized;
    },
    addSpec(state, {payload: {spec}}) {
      state.isDraftMode = true;
      state.specs = [...state.specs, spec];
    },
    updateSpec(state, {payload: {spec, selector}}) {
      state.isDraftMode = true;
      state.specs = state.specs.map(item => {
        if (item.selector === selector) {
          return {...spec, originalSelector: item.originalSelector};
        }
        return item;
      });
    },
    removeSpec(state, {payload: {selector}}) {
      state.isDraftMode = true;
      state.specs = state.specs.map(item => {
        if (item.selector === selector) {
          return {...item, isDraft: true, isDeleted: true};
        }
        return item;
      });
    },
    revertSpec(state, {payload: {originalSelector}}) {
      const initialSpec = state.initialSpecs.find(spec => spec.originalSelector === originalSelector);

      state.specs = initialSpec
        ? state.specs.map(spec => (spec.originalSelector === originalSelector ? initialSpec : spec))
        : state.specs.filter(spec => spec.selector !== originalSelector);

      const pendingChanges = state.specs.filter(({isDraft}) => isDraft).length;

      if (!pendingChanges) state.isDraftMode = false;
    },
    resetSpecs(state) {
      state.isDraftMode = false;
      state.specs = state.initialSpecs;
    },
    setAssertionResults(state, {payload}) {
      state.assertionResults = payload;
    },
    setSelectedSpec(state, {payload: assertionResult}) {
      if (assertionResult) state.selectedSpec = assertionResult.selector;
      else state.selectedSpec = undefined;
    },
    setSelectorSuggestions(state, {payload: selectorSuggestions}) {
      state.selectorSuggestions = selectorSuggestions;
    },
  },
  extraReducers: builder => {
    builder
      .addCase(TestSpecsActions.dryRun.fulfilled, (state, {payload}) => {
        state.assertionResults = payload;
      })
      .addCase(TestSpecsActions.publish.pending, state => {
        state.isDraftMode = false;
      })
      .addCase(TestSpecsActions.publish.fulfilled, (state, {payload: {result}}) => {
        const specs = assertionResultsToSpecs(result);

        state.assertionResults = result;
        state.initialSpecs = specs;
      })
      .addMatcher(
        ({type}) => type.startsWith('testDefinition') && type.endsWith('/pending'),
        state => {
          state.isLoading = true;
        }
      )
      .addMatcher(
        ({type}) => type.startsWith('testDefinition') && type.endsWith('/fulfilled'),
        state => {
          state.isLoading = false;
        }
      );
  },
});

export const {
  initSpecs,
  resetSpecs,
  addSpec,
  removeSpec,
  revertSpec,
  updateSpec,
  setAssertionResults,
  reset,
  setSelectedSpec,
  setIsInitialized,
  setSelectorSuggestions,
} = testSpecsSlice.actions;
export default testSpecsSlice.reducer;
