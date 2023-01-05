import {Model, TVariablesSchemas} from './Common.types';
import {TTest} from './Test.types';

export type TRawMissingVariables = TVariablesSchemas['MissingVariables'];
export type TMissingVariables = Model<TRawMissingVariables, {}>;

export type TRawVariables = TVariablesSchemas['Variables'];
export type TVariables = Model<
  TRawVariables,
  {
    missing: TMissingVariables[];
  }
>;

export type TRawTestVariables = TVariablesSchemas['TestVariables'];
export type TTestVariables = Model<
  TRawTestVariables,
  {
    variables: TVariables;
    test: TTest;
  }
>;
export type TTestVariablesMap = Record<string, TTestVariables>;

export type TTransactionVariables = {
  variables: TTestVariables[];
  hasMissingVariables: boolean;
};

export type TDraftVariables = {
  variables: Record<string, string>;
};
