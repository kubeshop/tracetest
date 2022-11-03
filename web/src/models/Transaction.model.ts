import {TRawTransaction, TTransaction} from 'types/Transaction.types';

function Transaction({id = '', name = '', description = '', version = 1, steps = []}: TRawTransaction): TTransaction {
  return {
    id,
    name,
    description,
    version,
    steps,
  };
}

export default Transaction;
