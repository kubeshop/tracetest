import {createSlice, PayloadAction} from '@reduxjs/toolkit';

import TestDefinitionActions, {TChange} from 'redux/actions/TestDefinition.actions';
import {TAssertionResults} from 'types/Assertion.types';
import {TTestDefinitionEntry} from 'types/TestDefinition.types';

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
    assertionList: resultList.flatMap(({assertion}) => [assertion]),
  }));
};

const testDefinitionSlice = createSlice({
  name: 'testDefinition',
  initialState,
  reducers: {
    reset() {
      return initialState;
    },
    initDefinitionList(state, {payload: {assertionResults}}: PayloadAction<{assertionResults: TAssertionResults}>) {
      const definitionList = assertionResultsToDefinitionList(assertionResults);

      state.initialDefinitionList = definitionList;
      state.definitionList = definitionList;
      state.isInitialized = true;
    },
    addDefinition(state, {payload: {definition}}: PayloadAction<{definition: TTestDefinitionEntry}>) {
      state.definitionList = [...state.definitionList, definition];
    },
    updateDefinition(
      state,
      {payload: {definition, selector}}: PayloadAction<{definition: TTestDefinitionEntry; selector: string}>
    ) {
      state.definitionList = state.definitionList.map(def => {
        if (def.selector === selector) return definition;

        return def;
      });
    },
    removeDefinition(state, {payload: {selector}}: PayloadAction<{selector: string}>) {
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
    setAssertionResults(state, {payload}: PayloadAction<TAssertionResults>) {
      state.assertionResults = payload;
    },
    clearAffectedSpans(state) {
      state.affectedSpans = [];
    },
    setAffectedSpans(state, {payload: spanIds}: PayloadAction<string[]>) {
      state.affectedSpans = spanIds;
    },
    setSelectedAssertion(state, {payload: selectorId}: PayloadAction<string>) {
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
  reset,
  clearAffectedSpans,
  setAffectedSpans,
  setSelectedAssertion,
} = testDefinitionSlice.actions;
export default testDefinitionSlice.reducer;
