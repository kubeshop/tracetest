import {TRawTransaction, TTransaction} from 'types/Transaction.types';
import TestSummary from './TestSummary.model';

function Transaction({
  id = '',
  name = '',
  description = '',
  version = 1,
  steps = [],
  createdAt = '',
  summary = {},
}: TRawTransaction): TTransaction {
  return {
    id,
    name,
    description,
    version,
    steps,
    createdAt,
    summary: TestSummary(summary),
  };
}

export default Transaction;
