import {TRawTransactionTestResult, TTransactionTestResult} from 'types/Transaction.types';

const TransactionRunResult = ({
  id = '',
  testId = '',
  result = 'running',
  name = '',
  version = 1,
  trigger,
}: TRawTransactionTestResult): TTransactionTestResult => ({
  id,
  testId,
  result,
  name,
  version,
  trigger: trigger!,
});

export default TransactionRunResult;
