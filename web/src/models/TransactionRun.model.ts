import {TRawTransactionRun, TTransactionRun} from 'types/TransactionRun.types';
import Environment from './Environment.model';
import TestRun from './TestRun.model';

const TransactionRunModel = ({
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
}: TRawTransactionRun): TTransactionRun => {
  return {
    id,
    createdAt,
    completedAt,
    state,
    steps: steps.map(step => TestRun(step)),
    environment: Environment(environment),
    metadata,
    version,
    pass,
    fail,
  };
};

export default TransactionRunModel;
