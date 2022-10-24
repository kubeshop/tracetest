import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';

import {TChange} from 'redux/actions/TestSpecs.actions';
import {TAssertionResultEntry, TAssertionResults} from './Assertion.types';
import {Model, TTestSchemas} from './Common.types';

export type TRawTestSpecs = TTestSchemas['TestSpecs'];

export type TTestSpecEntry = {
  assertions: string[];
  isDeleted?: boolean;
  isDraft: boolean;
  originalSelector?: string;
  selector: string;
};

export type TRawTestSpecEntry = {
  selector: {query: string};
  assertions: string[];
};

export type TTestSpecs = Model<TRawTestSpecs, {specs: TTestSpecEntry[]}>;

export interface ITestSpecsState {
  initialSpecs: TTestSpecEntry[];
  specs: TTestSpecEntry[];
  assertionResults?: TAssertionResults;
  changeList: TChange[];
  isLoading: boolean;
  isInitialized: boolean;
  selectedSpec: string | undefined;
  isDraftMode: boolean;
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
};
