import {TRawTransactionRun, TTransactionRun} from 'types/TransactionRun.types';
import Environment from './Environment.model';
import Test from './Test.model';
import TestRun from './TestRun.model';

const TransactionRunModel = ({
  id = '',
  createdAt = '',
  completedAt = '',
  state = 'CREATED',
  steps = [],
  stepRuns = [],
  environment = {},
  metadata = {},
  version = 1,
}: TRawTransactionRun): TTransactionRun => {
  return {
    id,
    createdAt,
    completedAt,
    state,
    steps: steps.map(step => Test(step)),
    stepRuns: stepRuns.map(stepRun => TestRun(stepRun)),
    environment: Environment(environment),
    metadata,
    version,
  };
};

export default TransactionRunModel;
