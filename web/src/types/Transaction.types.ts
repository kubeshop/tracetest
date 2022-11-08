import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {FormInstance} from 'antd';

import {Model, TTransactionsSchemas} from './Common.types';
import {TEnvironment, TRawEnvironment} from './Environment.types';
import {ICreateTestStep} from './Plugins.types';
import {TTrigger} from './Test.types';

export type TRawTransaction = TTransactionsSchemas['Transaction'];

export type TTransaction = Model<
  TRawTransaction,
  {
    id: string;
    name: string;
    description: string;
    version: number;
    createdAt?: string;
    steps: string[];
  }
>;

export type TRawTransactionTestResult = {
  id?: string;
  testId?: string;
  result?: 'success' | 'fail' | 'running';
  trigger?: TTrigger;
  name?: string;
  version?: number;
};

export type TTransactionTestResult = Model<TRawTransactionTestResult, {}>;

export type TRawTransactionRun = {
  id?: string;
  environment?: TRawEnvironment;
  results?: TTransactionTestResult[];
};

export type TTransactionRun = Model<
  TRawTransactionRun,
  {
    environment: TEnvironment;
  }
>;

export type TDraftTransaction = {
  steps?: string[];
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
