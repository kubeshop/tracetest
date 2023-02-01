import {RunErrorTypes} from 'types/TestRun.types';
import MissingVariables from './MissingVariables.model';

export type TRawRunError = any;
type RunError = {
  type: RunErrorTypes;
  missingVariables: MissingVariables;
};

const defaultError: RunError = {type: RunErrorTypes.Unknown, missingVariables: []};

const RunError = (error: TRawRunError): RunError => {
  if (!error) return defaultError;

  const {missingVariables} = error;

  return {
    ...defaultError,
    type: missingVariables ? RunErrorTypes.MissingVariables : RunErrorTypes.Unknown,
    missingVariables: MissingVariables(error),
  };
};

export default RunError;
