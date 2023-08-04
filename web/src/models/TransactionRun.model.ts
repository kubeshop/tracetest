import {Model, TTransactionsSchemas} from 'types/Common.types';
import VariableSet from './VariableSet.model';
import TestRun from './TestRun.model';

export type TRawTransactionResourceRun = TTransactionsSchemas['TransactionRun'];
type TransactionRun = Model<
  TRawTransactionResourceRun,
  {
    steps: TestRun[];
    variableSet?: VariableSet;
    metadata?: {[key: string]: string};
  }
>;

const TransactionRun = ({
  id = '',
  createdAt = '',
  completedAt = '',
  state = 'CREATED',
  steps = [],
  variableSet = {},
  metadata = {},
  version = 1,
  pass = 0,
  fail = 0,
  allStepsRequiredGatesPassed = false,
}: TRawTransactionResourceRun): TransactionRun => {
  return {
    id,
    createdAt,
    completedAt,
    state,
    steps: steps.map(step => TestRun(step)),
    variableSet: VariableSet.fromRun(variableSet),
    allStepsRequiredGatesPassed,
    metadata,
    version,
    pass,
    fail,
  };
};

export default TransactionRun;
