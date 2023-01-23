import {Model, TVariablesSchemas} from './Common.types';

export type TRawMissingVariables = TVariablesSchemas['MissingVariablesError'];
export type TMissingVariables = Model<TRawMissingVariables, {}>;

export type TRawMissingVariable = TVariablesSchemas['MissingVariable'];
export type TMissingVariable = Model<TRawMissingVariable, {}>;

export type TDraftVariables = {
  variables: Record<string, string>;
};
