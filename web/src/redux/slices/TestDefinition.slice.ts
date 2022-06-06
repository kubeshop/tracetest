import {createSlice} from '@reduxjs/toolkit';
import {TAssertionResults} from 'types/Assertion.types';
import {ITestDefinitionState, TTestDefinitionEntry, TTestDefinitionSliceActions} from 'types/TestDefinition.types';
import TestDefinitionActions from '../actions/TestDefinition.actions';

export const initialState: ITestDefinitionState = {
  initialDefinitionList: [],
  definitionList: [],
  changeList: [],
  isLoading: false,
  isInitialized: false,
  affectedSpans: [],
  selectedAssertion: '',
  selectedSpan: undefined,
  isDraftMode: false,
};

export const assertionResultsToDefinitionList = (assertionResults: TAssertionResults): TTestDefinitionEntry[] => {
  return assertionResults.resultList.map(({selector, resultList}) => ({
    isDraft: false,
    isDeleted: false,
    selector,
    originalSelector: selector,
    assertionList: resultList.flatMap(({assertion}) => [assertion]),
  }));
};

const testDefinitionSlice = createSlice<ITestDefinitionState, TTestDefinitionSliceActions, 'testDefinition'>({
  name: 'testDefinition',
  initialState,
  reducers: {
    reset() {
      return initialState;
    },
    initDefinitionList(state, {payload: {assertionResults}}) {
      const definitionList = assertionResultsToDefinitionList(assertionResults);

      state.initialDefinitionList = definitionList;
      state.definitionList = definitionList;
      state.isInitialized = true;
      state.isDraftMode = false;
    },
    addDefinition(state, {payload: {definition}}) {
      state.isDraftMode = true;
      state.definitionList = [...state.definitionList, definition];
    },
    updateDefinition(state, {payload: {definition, selector}}) {
      state.isDraftMode = true;
      state.definitionList = state.definitionList.map(def => {
        if (def.selector === selector)
          return {
            ...definition,
            originalSelector: def.originalSelector,
          };

        return def;
      });
    },
    removeDefinition(state, {payload: {selector}}) {
      state.isDraftMode = true;
      state.definitionList = state.definitionList.map(def => {
        if (def.selector === selector)
          return {
            ...def,
            isDraft: true,
            isDeleted: true,
          };

        return def;
      });
    },
    revertDefinition(state, {payload: {originalSelector}}) {
      const initialDefinition = state.initialDefinitionList.find(
        definition => definition.originalSelector === originalSelector
      );

      state.definitionList = initialDefinition
        ? state.definitionList.map(definition =>
            definition.originalSelector === originalSelector ? initialDefinition : definition
          )
        : state.definitionList.filter(definition => definition.selector !== originalSelector);

      const pendingChanges = state.definitionList.filter(({isDraft}) => isDraft).length;

      if (!pendingChanges) state.isDraftMode = false;
    },
    resetDefinitionList(state) {
      state.isDraftMode = false;
      state.definitionList = state.initialDefinitionList;
    },
    setAssertionResults(state, {payload}) {
      state.assertionResults = payload;
    },
    clearAffectedSpans(state) {
      state.affectedSpans = [];
    },
    setAffectedSpans(state, {payload: spanIds}) {
      state.affectedSpans = spanIds;
    },
    setSelectedAssertion(state, {payload: selectorId}) {
      const assertionResult = state?.assertionResults?.resultList?.find(assertion => assertion.selector === selectorId);
      const spanIds = assertionResult?.spanIds ?? [];
      state.selectedAssertion = selectorId;
      state.affectedSpans = spanIds;
    },
    setSelectedSpan(state, {payload: span}) {
      state.selectedSpan = span;
    },
  },
  extraReducers: builder => {
    builder
      .addCase(TestDefinitionActions.dryRun.fulfilled, (state, {payload}) => {
        state.assertionResults = payload;
      })
      .addCase(TestDefinitionActions.publish.pending, state => {
        state.isDraftMode = false;
      })
      .addCase(TestDefinitionActions.publish.fulfilled, (state, {payload: {result}}) => {
        const definitionList = assertionResultsToDefinitionList(result);

        state.assertionResults = result;
        state.initialDefinitionList = definitionList;
      })
      .addMatcher(
        action => action.type.startsWith('testDefinition') && action.type.endsWith('/pending'),
        state => {
          state.isLoading = true;
        }
      )
      .addMatcher(
        action => action.type.startsWith('testDefinition') && action.type.endsWith('/fulfilled'),
        state => {
          state.isLoading = false;
        }
      );
  },
});

export const {
  addDefinition,
  removeDefinition,
  updateDefinition,
  setAssertionResults,
  initDefinitionList,
  resetDefinitionList,
  revertDefinition,
  reset,
  clearAffectedSpans,
  setAffectedSpans,
  setSelectedAssertion,
  setSelectedSpan,
} = testDefinitionSlice.actions;
export default testDefinitionSlice.reducer;
