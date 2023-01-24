import {Model, TVariablesSchemas} from './Common.types';
import {TTest} from './Test.types';

export type TRawMissingVariables = TVariablesSchemas['MissingVariablesError'];
export type TMissingVariables = Model<TRawMissingVariables, {
  missingVariables: TMissingVariable[];
}>;

export type TRawVariable = TVariablesSchemas['Variable'];
export type TVariable = Model<TRawVariable, {}>;

export type TRawMissingVariable = TVariablesSchemas['MissingVariable'];
export type TMissingVariable = Model<
  TRawMissingVariable,
  {
    variables: TVariable[];
  }
>;

export type TTestVariables = {
  variables: TVariable[];
  test: TTest;
};
export type TTestVariablesMap = Record<string, TTestVariables>;

export type TDraftVariables = {
  variables: Record<string, string>;
};
