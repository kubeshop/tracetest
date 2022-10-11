import {TRawTransaction, TTransaction} from 'types/Transaction.types';

const Transaction = ({id = '', name = '', description = '', version = 1}: TRawTransaction): TTransaction => ({
  id,
  name,
  description,
  version,
});

export default Transaction;
