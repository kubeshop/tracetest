import Test from 'models/Test.model';
import {Model, TVariablesSchemas} from './Common.types';

export type TVariable = Model<TRawVariable, {}>;
export type TRawVariable = TVariablesSchemas['Variable'];

export type TTestVariables = {
  variables: TVariable[];
  test: Test;
};
export type TTestVariablesMap = Record<string, TTestVariables>;

export type TDraftVariables = {
  variables: Record<string, string>;
};
