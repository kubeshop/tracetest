import {TRawTestVariables, TTransactionVariables} from '../types/Variables.types';
import TestVariables from './TestVariables.model';

const TransactionVariables = (testVariables: TRawTestVariables[] = []): TTransactionVariables => {
  const variables = testVariables.map(testVariable => TestVariables(testVariable));
  const hasMissingVariables = !!variables.find(({variables: {missing}}) => missing.length > 0);

  return {
    hasMissingVariables,
    variables,
  };
};

export default TransactionVariables;
