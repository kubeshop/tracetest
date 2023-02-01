import {Model, TVariablesSchemas} from 'types/Common.types';
import {TVariable} from 'types/Variables.types';

export type TRawMissingVariables = TVariablesSchemas['MissingVariablesError'];
export type TRawMissingVariable = TVariablesSchemas['MissingVariable'];
export type MissingVariable = Model<
  TRawMissingVariable,
  {
    variables: TVariable[];
  }
>;
type MissingVariables = MissingVariable[];

const MissingVariables = ({missingVariables = []}: TRawMissingVariables = {}): MissingVariables => {
  return missingVariables.map(({testId = '', variables = []}) => ({
    testId,
    variables: variables.map(({key = '', defaultValue = ''}) => ({key, defaultValue})),
  }));
};

export default MissingVariables;
