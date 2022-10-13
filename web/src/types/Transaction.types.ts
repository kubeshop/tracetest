import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {FormInstance} from 'antd';
import {Model} from './Common.types';
import {ICreateTestStep} from './Plugins.types';
import {TTest} from './Test.types';

export type TRawTransaction = {
  id?: string;
  name?: string;
  description?: string;
  version?: number;
};
export interface TTransaction extends Model<TRawTransaction, ITransaction> {}

interface TransactionStep extends TTest {
  result: 'success' | 'fail' | 'running';
}

export interface ITransaction {
  id: string;
  name: string;
  description: string;
  version: number;
  steps: TransactionStep[];
  env: Record<string, string>;
}

export type TDraftTransaction = {
  tests?: string[];
  name?: string;
  description?: string;
};

export interface ICreateTransactionState {
  draftTransaction: TDraftTransaction;
  stepNumber: number;
  isFormValid: boolean;
  stepList: ICreateTestStep[];
}

export interface IBasicValues {
  name: string;
  description: string;
  tests: string[];
}

export type TDraftTransactionForm = FormInstance<TDraftTransaction>;

export type TCreateTransactionSliceActions = {
  reset: CaseReducer<ICreateTransactionState>;
  setStepNumber: CaseReducer<ICreateTransactionState, PayloadAction<{stepNumber: number; completeStep?: boolean}>>;
  setDraftTransaction: CaseReducer<ICreateTransactionState, PayloadAction<{draftTransaction: TDraftTransaction}>>;
  setIsFormValid: CaseReducer<ICreateTransactionState, PayloadAction<{isValid: boolean}>>;
};
