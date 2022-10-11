import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {FormInstance} from 'antd';
import {ICreateTestStep} from './Plugins.types';
import {Model} from './Common.types';

export type TRawTransaction = {
  id?: string;
  name?: string;
  description?: string;
  version?: number;
};
export type TTransaction = Model<TRawTransaction, {}>;

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
