import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';

import {TChange} from 'redux/actions/TestDefinition.actions';
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

export interface ITestDefinitionState {
  initialDefinitionList: TTestSpecEntry[];
  definitionList: TTestSpecEntry[];
  assertionResults?: TAssertionResults;
  changeList: TChange[];
  isLoading: boolean;
  isInitialized: boolean;
  selectedAssertion: string | undefined;
  isDraftMode: boolean;
}

export type TTestDefinitionSliceActions = {
  reset: CaseReducer<ITestDefinitionState>;
  initDefinitionList: CaseReducer<ITestDefinitionState, PayloadAction<{assertionResults: TAssertionResults}>>;
  addDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{definition: TTestSpecEntry}>>;
  updateDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{definition: TTestSpecEntry; selector: string}>>;
  removeDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{selector: string}>>;
  revertDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{originalSelector: string}>>;
  resetDefinitionList: CaseReducer<ITestDefinitionState>;
  setAssertionResults: CaseReducer<ITestDefinitionState, PayloadAction<TAssertionResults>>;
  setSelectedAssertion: CaseReducer<ITestDefinitionState, PayloadAction<TAssertionResultEntry | undefined>>;
  setIsInitialized: CaseReducer<ITestDefinitionState, PayloadAction<{isInitialized: boolean}>>;
};
