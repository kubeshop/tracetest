import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {FormInstance} from 'antd';

import {Model, TTransactionsSchemas} from './Common.types';
import {ICreateTestStep} from './Plugins.types';
import {TTest} from './Test.types';

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
    // steps: TransactionStep[]; // TODO define if this should be part of Transaction Run
    env?: Record<string, string>; // TODO define if this should be part of Transaction Run
  }
>;

interface TransactionStep extends TTest {
  result: 'success' | 'fail' | 'running';
}

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
