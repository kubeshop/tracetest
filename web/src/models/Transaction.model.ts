import {Model, TTransactionsSchemas} from 'types/Common.types';
import Summary from './Summary.model';
import Test from './Test.model';

export type TRawTransaction = TTransactionsSchemas['Transaction'];
type Transaction = Model<
  TRawTransaction,
  {
    steps: Test[];
    summary: Summary;
  }
>;

function Transaction({
  id = '',
  name = '',
  description = '',
  version = 1,
  steps = [],
  createdAt = '',
  summary = {},
}: TRawTransaction): Transaction {
  return {
    id,
    name,
    description,
    version,
    steps: steps.map(step => Test(step)),
    createdAt,
    summary: Summary(summary),
  };
}

export default Transaction;
