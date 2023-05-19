import {Model, TTransactionsSchemas} from 'types/Common.types';
import Summary from './Summary.model';
import Test from './Test.model';

export type TRawTransactionResource = TTransactionsSchemas['TransactionResource'];
export type TRawTransaction = TTransactionsSchemas['Transaction'];
type Transaction = Model<
  TTransactionsSchemas['Transaction'],
  {
    steps: string[];
    fullSteps: Test[];
    summary: Summary;
  }
>;

function Transaction({
  spec: rawTransaction = {},
}: TRawTransactionResource): Transaction {
  return Transaction.FromRawTransaction(rawTransaction);
}

Transaction.FromRawTransaction = ({
  id = '',
  name = '',
  description = '',
  version = 1,
  createdAt = '',
  summary = {},
  steps = [],
  fullSteps = [],
}: TRawTransaction): Transaction => {
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
};

export default Transaction;
