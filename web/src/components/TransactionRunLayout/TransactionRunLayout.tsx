import {useMemo} from 'react';
import TransactionHeader from 'components/TransactionHeader';
import useDocumentTitle from 'hooks/useDocumentTitle';
import TransactionService from 'services/Transaction.service';
import Transaction from 'models/Transaction.model';
import TransactionRun from 'models/TransactionRun.model';
import * as S from './TransactionRunLayout.styled';
import EditTransaction from '../EditTransaction';
import TransactionRunResult from '../TransactionRunResult/TransactionRunResult';

interface IProps {
  transaction: Transaction;
  transactionRun: TransactionRun;
}

const TransactionRunDetailLayout = ({transaction, transactionRun}: IProps) => {
  useDocumentTitle(`${transaction.name} - ${transactionRun.state}`);
  const draftTransaction = useMemo(() => TransactionService.getInitialValues(transaction), [transaction]);

  return (
    <>
      <TransactionHeader transactionRun={transactionRun} transaction={transaction} />
      <S.Wrapper>
        <S.Container>
          <S.SectionLeft>
            <EditTransaction transaction={draftTransaction} transactionRun={transactionRun} />
          </S.SectionLeft>
          <S.SectionRight>
            <TransactionRunResult transaction={transaction} transactionRun={transactionRun} />
          </S.SectionRight>
        </S.Container>
      </S.Wrapper>
    </>
  );
};

export default TransactionRunDetailLayout;
