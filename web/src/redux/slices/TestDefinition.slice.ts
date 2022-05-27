import {createSlice, PayloadAction} from '@reduxjs/toolkit';
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
}

export const initialState: ITestDefinitionState = {
  initialDefinitionList: [],
  definitionList: [],
  changeList: [],
  isLoading: false,
  isInitialized: false,
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
} = testDefinitionSlice.actions;
export default testDefinitionSlice.reducer;
