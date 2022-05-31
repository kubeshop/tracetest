import {CaseReducer, createSlice, PayloadAction} from '@reduxjs/toolkit';

import {TAssertionResults} from '../../types/Assertion.types';
import {TTestDefinitionEntry} from '../../types/TestDefinition.types';
import TestDefinitionActions, {TChange} from '../actions/TestDefinition.actions';

interface ITestDefinitionState {
  initialDefinitionList: TTestDefinitionEntry[];
  definitionList: TTestDefinitionEntry[];
  assertionResults?: TAssertionResults;
  changeList: TChange[];
  isLoading: boolean;
  isInitialized: boolean;
  affectedSpans: string[];
  selectedAssertion: string;
}

export const initialState: ITestDefinitionState = {
  initialDefinitionList: [],
  definitionList: [],
  changeList: [],
  isLoading: false,
  isInitialized: false,
  affectedSpans: [],
  selectedAssertion: '',
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

const testDefinitionSlice = createSlice<
  ITestDefinitionState,
  {
    reset: CaseReducer<ITestDefinitionState>;
    initDefinitionList: CaseReducer<ITestDefinitionState, PayloadAction<{assertionResults: TAssertionResults}>>;
    addDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{definition: TTestDefinitionEntry}>>;
    updateDefinition: CaseReducer<
      ITestDefinitionState,
      PayloadAction<{definition: TTestDefinitionEntry; selector: string}>
    >;
    removeDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{selector: string}>>;
    revertDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{originalSelector: string}>>;
    resetDefinitionList: CaseReducer<ITestDefinitionState>;
    setAssertionResults: CaseReducer<ITestDefinitionState, PayloadAction<TAssertionResults>>;
    clearAffectedSpans: CaseReducer<ITestDefinitionState>;
    setAffectedSpans: CaseReducer<ITestDefinitionState, PayloadAction<string[]>>;
    setSelectedAssertion: CaseReducer<ITestDefinitionState, PayloadAction<string>>;
  },
  'testDefinition'
>({
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
    },
    addDefinition(state, {payload: {definition}}) {
      state.definitionList = [...state.definitionList, definition];
    },
    revertDefinition(state, {payload: {originalSelector}}) {
      const initialDefinition = state.initialDefinitionList.find(
        definition => definition.originalSelector === originalSelector
      );

      state.definitionList = initialDefinition
        ? state.definitionList.map(definition =>
            definition.originalSelector === originalSelector ? initialDefinition : definition
          )
        : state.definitionList.filter(definition => definition.originalSelector === originalSelector);
    },
    updateDefinition(state, {payload: {definition, selector}}) {
      state.definitionList = state.definitionList.map(def => {
        if (def.originalSelector === selector)
          return {
            ...definition,
            originalSelector: def.originalSelector,
          };

        return def;
      });
    },
    removeDefinition(state, {payload: {selector}}) {
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
    resetDefinitionList(state) {
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
  },
  extraReducers: builder => {
    builder
      .addCase(TestDefinitionActions.dryRun.fulfilled, (state, {payload}) => {
        state.assertionResults = payload;
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
} = testDefinitionSlice.actions;
export default testDefinitionSlice.reducer;
