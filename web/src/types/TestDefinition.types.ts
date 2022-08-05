import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {TChange} from '../redux/actions/TestDefinition.actions';
import {TAssertion, TAssertionResultEntry, TAssertionResults} from './Assertion.types';
import {Model, TTestSchemas} from './Common.types';

export type TRawTestDefinition = TTestSchemas['TestDefinition'];

export type TTestDefinitionEntry = {
  originalSelector?: string;
  selector: string;
  assertionList: TAssertion[];
  isDraft: boolean;
  isDeleted?: boolean;
};

export type TRawTestDefinitionEntry = {
  selector: {query: string};
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
  selectedAssertion: string | undefined;
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
  revertDefinition: CaseReducer<ITestDefinitionState, PayloadAction<{originalSelector: string}>>;
  resetDefinitionList: CaseReducer<ITestDefinitionState>;
  setAssertionResults: CaseReducer<ITestDefinitionState, PayloadAction<TAssertionResults>>;
  setSelectedAssertion: CaseReducer<ITestDefinitionState, PayloadAction<TAssertionResultEntry | undefined>>;
};
