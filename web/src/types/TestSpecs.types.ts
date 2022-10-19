import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';

import {TChange} from 'redux/actions/TestSpecs.actions';
import {TAssertion, TAssertionResultEntry, TAssertionResults} from './Assertion.types';
import {Model, TTestSchemas} from './Common.types';

export type TRawTestSpecs = TTestSchemas['TestSpecs'];

export type TTestSpecEntry = {
  assertions: TAssertion[];
  isDeleted?: boolean;
  isDraft: boolean;
  originalSelector?: string;
  selector: string;
};

export type TRawTestSpecEntry = {
  selector: {query: string};
  assertions: TAssertion[];
};

export type TTestSpecs = Model<TRawTestSpecs, {specs: TTestSpecEntry[]}>;

export interface ISuggestion {
  query: string;
  title: string;
}

export interface ITestSpecsState {
  initialSpecs: TTestSpecEntry[];
  specs: TTestSpecEntry[];
  assertionResults?: TAssertionResults;
  changeList: TChange[];
  isLoading: boolean;
  isInitialized: boolean;
  selectedSpec: string | undefined;
  isDraftMode: boolean;
  selectorSuggestions: ISuggestion[];
}

export type TTestSpecsSliceActions = {
  reset: CaseReducer<ITestSpecsState>;
  initSpecs: CaseReducer<ITestSpecsState, PayloadAction<{assertionResults: TAssertionResults}>>;
  addSpec: CaseReducer<ITestSpecsState, PayloadAction<{spec: TTestSpecEntry}>>;
  updateSpec: CaseReducer<ITestSpecsState, PayloadAction<{spec: TTestSpecEntry; selector: string}>>;
  removeSpec: CaseReducer<ITestSpecsState, PayloadAction<{selector: string}>>;
  revertSpec: CaseReducer<ITestSpecsState, PayloadAction<{originalSelector: string}>>;
  resetSpecs: CaseReducer<ITestSpecsState>;
  setAssertionResults: CaseReducer<ITestSpecsState, PayloadAction<TAssertionResults>>;
  setSelectedSpec: CaseReducer<ITestSpecsState, PayloadAction<TAssertionResultEntry | undefined>>;
  setIsInitialized: CaseReducer<ITestSpecsState, PayloadAction<{isInitialized: boolean}>>;
  setSelectorSuggestions: CaseReducer<ITestSpecsState, PayloadAction<ISuggestion[]>>;
};
