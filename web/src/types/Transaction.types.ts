import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {FormInstance} from 'antd';

import {Model, TTransactionsSchemas} from './Common.types';
import {ICreateTestStep} from './Plugins.types';
import { TTest, TSummary } from './Test.types';

export type TRawTransaction = TTransactionsSchemas['Transaction'];

export type TTransaction = Model<TRawTransaction, {
  steps: TTest[];
  summary: TSummary;
}>;

export type TDraftTransaction = {
  steps?: TTest[];
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
