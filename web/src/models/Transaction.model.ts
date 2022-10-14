import {TTransaction} from 'types/Transaction.types';

const Transaction = ({id = '', name = '', description = '', version = 1, ...data}: TTransaction): TTransaction => ({
  id,
  name,
  description,
  version,
  ...data,
});

export default Transaction;
