import {useNavigate} from 'react-router-dom';
import TransactionHeader from 'components/TransactionHeader';
import {TTransaction} from 'types/Transaction.types';
import {TTransactionRun} from 'types/TransactionRun.types';
import * as S from './TransactionRunLayout.styled';
import EditTransaction from '../EditTransaction';
import TransactionRunResult from '../TransactionRunResult/TransactionRunResult';

interface IProps {
  transaction: TTransaction;
  transactionRun: TTransactionRun;
}

const TransactionRunDetailLayout = ({transaction, transaction: {id: transactionId}, transactionRun}: IProps) => {
  const navigate = useNavigate();

  return (
    <>
      <TransactionHeader onBack={() => navigate(`/transaction/${transactionId}`)} />
      <S.Wrapper>
        <S.Container>
          <S.SectionLeft>
            <EditTransaction transaction={transaction} />
          </S.SectionLeft>
          <S.SectionRight>
            <TransactionRunResult transactionRun={transactionRun} />
          </S.SectionRight>
        </S.Container>
      </S.Wrapper>
    </>
  );
};

export default TransactionRunDetailLayout;
