import {useMemo} from 'react';
import EditTransaction from 'components/EditTransaction';
import TransactionRunResult from 'components/TransactionRunResult';
import useDocumentTitle from 'hooks/useDocumentTitle';
import {useTransaction} from 'providers/Transaction/Transaction.provider';
import {useTransactionRun} from 'providers/TransactionRun/TransactionRun.provider';
import TransactionService from 'services/Transaction.service';
import * as S from './TransactionRunOverview.styled';

const Content = () => {
  const {transaction} = useTransaction();
  const {transactionRun} = useTransactionRun();
  useDocumentTitle(`${transaction.name} - ${transactionRun.state}`);
  const draftTransaction = useMemo(() => TransactionService.getInitialValues(transaction), [transaction]);

  return (
    <S.Container>
      <S.SectionLeft>
        <EditTransaction transaction={draftTransaction} transactionRun={transactionRun} />
      </S.SectionLeft>
      <S.SectionRight>
        <TransactionRunResult transaction={transaction} transactionRun={transactionRun} />
      </S.SectionRight>
    </S.Container>
  );
};

export default Content;
