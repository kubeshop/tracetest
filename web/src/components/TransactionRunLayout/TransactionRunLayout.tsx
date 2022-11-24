import {useMemo} from 'react';
import {useNavigate} from 'react-router-dom';
import TransactionHeader from 'components/TransactionHeader';
import {TTransaction} from 'types/Transaction.types';
import {TTransactionRun} from 'types/TransactionRun.types';
import useDocumentTitle from 'hooks/useDocumentTitle';
import TransactionService from 'services/Transaction.service';
import * as S from './TransactionRunLayout.styled';
import EditTransaction from '../EditTransaction';
import TransactionRunResult from '../TransactionRunResult/TransactionRunResult';

interface IProps {
  transaction: TTransaction;
  transactionRun: TTransactionRun;
}

const TransactionRunDetailLayout = ({transaction, transaction: {id: transactionId}, transactionRun}: IProps) => {
  const navigate = useNavigate();
  useDocumentTitle(`${transaction.name} - ${transactionRun.state}`);
  const draftTransaction = useMemo(() => TransactionService.getInitialValues(transaction), [transaction]);

  return (
    <>
      <TransactionHeader
        onBack={() => navigate(`/transaction/${transactionId}`)}
        transactionRun={transactionRun}
        transaction={transaction}
      />
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
