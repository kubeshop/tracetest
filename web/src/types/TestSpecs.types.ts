import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';

import {TChange} from 'redux/actions/TestSpecs.actions';
import TestSpecs, {TTestSpecEntry} from 'models/TestSpecs.model';
import AssertionResults, {TAssertionResultEntry} from 'models/AssertionResults.model';

export interface ISuggestion {
  query: string;
  title: string;
}

export interface ITestSpecsState {
  initialSpecs: TTestSpecEntry[];
  specs: TTestSpecEntry[];
  assertionResults?: AssertionResults;
  changeList: TChange[];
  isLoading: boolean;
  isInitialized: boolean;
  selectedSpec: string | undefined;
  isDraftMode: boolean;
  selectorSuggestions: ISuggestion[];
  prevSelector: string;
}

export type TTestSpecsSliceActions = {
  reset: CaseReducer<ITestSpecsState>;
  initSpecs: CaseReducer<ITestSpecsState, PayloadAction<{assertionResults: AssertionResults; specs: TestSpecs}>>;
  addSpec: CaseReducer<ITestSpecsState, PayloadAction<{spec: TTestSpecEntry}>>;
  updateSpec: CaseReducer<ITestSpecsState, PayloadAction<{spec: TTestSpecEntry; selector: string}>>;
  removeSpec: CaseReducer<ITestSpecsState, PayloadAction<{selector: string}>>;
  revertSpec: CaseReducer<ITestSpecsState, PayloadAction<{originalSelector: string}>>;
  resetSpecs: CaseReducer<ITestSpecsState>;
  setAssertionResults: CaseReducer<ITestSpecsState, PayloadAction<AssertionResults>>;
  setSelectedSpec: CaseReducer<ITestSpecsState, PayloadAction<TAssertionResultEntry | undefined>>;
  setIsInitialized: CaseReducer<ITestSpecsState, PayloadAction<{isInitialized: boolean}>>;
  setSelectorSuggestions: CaseReducer<ITestSpecsState, PayloadAction<ISuggestion[]>>;
  setPrevSelector: CaseReducer<ITestSpecsState, PayloadAction<{prevSelector: string}>>;
};
