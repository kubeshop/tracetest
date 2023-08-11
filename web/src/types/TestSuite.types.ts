import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {FormInstance} from 'antd';
import {ICreateTestStep} from './Plugins.types';

export type TDraftTestSuite = {
  steps?: string[];
  name?: string;
  description?: string;
};

export interface ICreateTestSuiteState {
  draft: TDraftTestSuite;
  stepNumber: number;
  isFormValid: boolean;
  stepList: ICreateTestStep[];
}

export interface IBasicValues {
  name: string;
  description: string;
  tests: string[];
}

export type TDraftTestSuiteForm = FormInstance<TDraftTestSuite>;

export type TCreateTestSuiteSliceActions = {
  reset: CaseReducer<ICreateTestSuiteState>;
  setStepNumber: CaseReducer<ICreateTestSuiteState, PayloadAction<{stepNumber: number; completeStep?: boolean}>>;
  setDraft: CaseReducer<ICreateTestSuiteState, PayloadAction<{draft: TDraftTestSuite}>>;
  setIsFormValid: CaseReducer<ICreateTestSuiteState, PayloadAction<{isValid: boolean}>>;
};
