import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {TChange} from '../redux/actions/TestDefinition.actions';
import {TAssertion, TAssertionResults} from './Assertion.types';
import {Model, TTestSchemas} from './Common.types';
import {TSpan} from './Span.types';

export type TRawTestDefinition = TTestSchemas['TestDefinition'];

export type TTestDefinitionEntry = {
  originalSelector?: string;
  selector: string;
  assertionList: TAssertion[];
  isDraft: boolean;
  isDeleted?: boolean;
};

export type TRawTestDefinitionEntry = {
  selector: string;
  assertions: TAssertion[];
};

export type TTestDefinition = Model<
  TRawTestDefinition,
  {
    definitionList: TTestDefinitionEntry[];
    definitions?: TRawTestDefinition;
  }
>;

export interface ITestDefinitionState {
  initialDefinitionList: TTestDefinitionEntry[];
  definitionList: TTestDefinitionEntry[];
  assertionResults?: TAssertionResults;
  changeList: TChange[];
  isLoading: boolean;
  isInitialized: boolean;
  affectedSpans: string[];
  selectedAssertion: string;
  selectedSpan?: TSpan;
  isDraftMode: boolean;
}

export type TTestDefinitionSliceActions = {
  reset: CaseReducer<ITestDefinitionState>;
  initDefinitionList: CaseReducer<ITestDefinitionState, PayloadAction<{assertionResults: TAssertionResults}>>;
  addDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{definition: TTestDefinitionEntry}>>;
  updateDefinition: CaseReducer<
    ITestDefinitionState,
    PayloadAction<{definition: TTestDefinitionEntry; selector: string}>
  >;
  removeDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{selector: string}>>;
  revertDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{originalSelector: string; selector: string}>>;
  resetDefinitionList: CaseReducer<ITestDefinitionState>;
  setAssertionResults: CaseReducer<ITestDefinitionState, PayloadAction<TAssertionResults>>;
  clearAffectedSpans: CaseReducer<ITestDefinitionState>;
  setAffectedSpans: CaseReducer<ITestDefinitionState, PayloadAction<string[]>>;
  setSelectedAssertion: CaseReducer<ITestDefinitionState, PayloadAction<string>>;
  setSelectedSpan: CaseReducer<ITestDefinitionState, PayloadAction<TSpan | undefined>>;
};
