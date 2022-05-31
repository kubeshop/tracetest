import {CaseReducer, createSlice, PayloadAction} from '@reduxjs/toolkit';
import getUuidByString from 'uuid-by-string';
import {TAssertionResults} from '../../types/Assertion.types';
import {TTestDefinitionEntry} from '../../types/TestDefinition.types';
import TestDefinitionActions, {TChange} from '../actions/TestDefinition.actions';

export interface ITestDefinitionState {
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
    id: getUuidByString(selector),
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
    revertDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{id?: string}>>;
    resetDefinitionList: CaseReducer<ITestDefinitionState>;
    setAssertionResults: CaseReducer<ITestDefinitionState, PayloadAction<TAssertionResults>>;
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
    updateDefinition(state, {payload: {definition, selector}}) {
      state.definitionList = state.definitionList.map(def => {
        if (def.selector === selector)
          return {
            ...def,
            ...definition,
          };

        return def;
      });
    },
    revertDefinition(state, {payload: {id}}) {
      state.definitionList = state.definitionList.map(d => {
        return d?.id === id ? state.initialDefinitionList.find(f => f.id === id) || d : d;
      });
    },
    removeDefinition(state: ITestDefinitionState, {payload: {selector}}: PayloadAction<{selector: string}>) {
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
      console.log('here');
      state.assertionResults = payload;
    },
  },
  extraReducers: builder => {
    builder
      .addCase(TestDefinitionActions.dryRun.fulfilled, (state: ITestDefinitionState, {payload}) => {
        state.assertionResults = payload;
      })
      .addCase(TestDefinitionActions.publish.fulfilled, (state: ITestDefinitionState, {payload: {result}}) => {
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
  revertDefinition,
  updateDefinition,
  setAssertionResults,
  initDefinitionList,
  resetDefinitionList,
  reset,
} = testDefinitionSlice.actions;
export default testDefinitionSlice.reducer;
