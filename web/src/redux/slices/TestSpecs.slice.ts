import {createSlice} from '@reduxjs/toolkit';

import AssertionResults from 'models/AssertionResults.model';
import TestSpecs, {TTestSpecEntry} from 'models/TestSpecs.model';
import {TTestSpecsSliceActions, ITestSpecsState} from 'types/TestSpecs.types';
import TestSpecsActions from '../actions/TestSpecs.actions';

export const initialState: ITestSpecsState = {
  initialSpecs: [],
  specs: [],
  changeList: [],
  isLoading: false,
  isInitialized: false,
  selectedSpec: undefined,
  selectedSpanResult: undefined,
  isDraftMode: false,
  selectorSuggestions: [],
  prevSelector: '',
};

export const assertionResultsToSpecs = (assertionResults: AssertionResults, specs: TestSpecs): TTestSpecEntry[] => {
  return assertionResults.resultList.map(({selector, resultList}, index) => ({
    isDraft: false,
    isDeleted: false,
    selector,
    originalSelector: selector,
    assertions: resultList.flatMap(({assertion}) => [assertion]),
    name: specs?.specs?.[index]?.name ?? '',
  }));
};

const testSpecsSlice = createSlice<ITestSpecsState, TTestSpecsSliceActions, 'testSpecs'>({
  name: 'testSpecs',
  initialState,
  reducers: {
    reset() {
      return initialState;
    },
    initSpecs(state, {payload: {assertionResults, specs}}) {
      const specsFromResults = assertionResultsToSpecs(assertionResults, specs);

      state.initialSpecs = specsFromResults;
      state.specs = specsFromResults;
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
    setSelectedSpanResult(state, {payload: spanResult}) {
      if (spanResult) state.selectedSpanResult = spanResult;
      else state.selectedSpanResult = undefined;
    },
    setSelectorSuggestions(state, {payload: selectorSuggestions}) {
      state.selectorSuggestions = selectorSuggestions;
    },
    setPrevSelector(state, {payload: {prevSelector}}) {
      state.prevSelector = prevSelector;
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
        const specs = assertionResultsToSpecs(result, TestSpecs({}));

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
  setPrevSelector,
  setSelectedSpanResult,
} = testSpecsSlice.actions;
export default testSpecsSlice.reducer;
