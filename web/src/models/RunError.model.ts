import {RunErrorTypes, TRawRunError, TRunError} from 'types/TestRun.types';
import {TMissingVariable, TRawMissingVariables} from 'types/Variables.types';

export const MissingVariables = ({missingVariables = []}: TRawMissingVariables = {}): TMissingVariable[] => {
  return missingVariables.map(({testId = '', variables = []}) => ({
    testId,
    variables: variables.map(({key = '', defaultValue = ''}) => ({key, defaultValue})),
  }));
};

const defaultError: TRunError = {type: RunErrorTypes.Unknown, missingVariables: []};

const RunError = (error: TRawRunError): TRunError => {
  if (!error) return defaultError;

  const {missingVariables} = error;

  return {
    ...defaultError,
    type: missingVariables ? RunErrorTypes.MissingVariables : RunErrorTypes.Unknown,
    missingVariables: MissingVariables(error),
  };
};

export default RunError;
