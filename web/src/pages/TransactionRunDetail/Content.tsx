import TransactionRunLayout from 'components/TransactionRunLayout';
import {useTransaction} from 'providers/Transaction/Transaction.provider';
import {useTransactionRun} from 'providers/TransactionRun/TransactionRun.provider';

const Content = () => {
  const {transaction} = useTransaction();
  const {transactionRun} = useTransactionRun();

  return transaction ? <TransactionRunLayout transaction={transaction} transactionRun={transactionRun!} /> : null;
};

export default Content;
