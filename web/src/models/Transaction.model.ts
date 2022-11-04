import {TRawTransaction, TTransaction} from 'types/Transaction.types';

function Transaction({
  id = '',
  name = '',
  description = '',
  version = 1,
  steps = [],
  createdAt = '',
}: TRawTransaction): TTransaction {
  return {
    id,
    name,
    description,
    version,
    steps,
    createdAt,
  };
}

export default Transaction;
