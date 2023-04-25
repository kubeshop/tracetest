import {Model, TTransactionsSchemas} from 'types/Common.types';
import Environment from './Environment.model';
import TestRun from './TestRun.model';

export type TRawTransactionRun = TTransactionsSchemas['TransactionRun'];
type TransactionRun = Model<
  TRawTransactionRun,
  {
    steps: TestRun[];
    environment?: Environment;
    metadata?: {[key: string]: string};
  }
>;

const TransactionRun = ({
  id = '',
  createdAt = '',
  completedAt = '',
  state = 'CREATED',
  steps = [],
  environment = {},
  metadata = {},
  version = 1,
  pass = 0,
  fail = 0,
}: TRawTransactionRun): TransactionRun => {
  return {
    id,
    createdAt,
    completedAt,
    state,
    steps: steps.map(step => TestRun(step)),
    environment: Environment.fromRun(environment),
    metadata,
    version,
    pass,
    fail,
  };
};

export default TransactionRun;
