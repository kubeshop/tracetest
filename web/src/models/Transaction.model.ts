import {TRawTransaction, TTransaction} from 'types/Transaction.types';
import TestSummary from './TestSummary.model';
import Test from './Test.model';

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
    steps: steps.map(step => Test(step)),
    createdAt,
    summary: TestSummary(summary),
  };
}

export default Transaction;
