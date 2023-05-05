import {Model, TTransactionsSchemas} from 'types/Common.types';
import Summary from './Summary.model';
import Test from './Test.model';

export type TRawTransaction = TTransactionsSchemas['TransactionResource'];
type Transaction = Model<
  TTransactionsSchemas['Transaction'],
  {
    steps: string[];
    fullSteps: Test[];
    summary: Summary;
  }
>;

function Transaction({
  spec: {
    id = '',
    name = '',
    description = '',
    version = 1,
    steps = [],
    fullSteps = [],
    createdAt = '',
    summary = {},
  } = {},
}: TRawTransaction): Transaction {
  return {
    id,
    name,
    description,
    version,
    steps,
    fullSteps: fullSteps.map(step => Test(step)),
    createdAt,
    summary: Summary(summary),
  };
}

export default Transaction;
