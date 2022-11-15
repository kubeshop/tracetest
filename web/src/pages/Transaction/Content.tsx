import {Button} from 'antd';
import {useCallback, useMemo} from 'react';
import {useNavigate} from 'react-router-dom';

import PaginatedList from 'components/PaginatedList';
import TransactionRunCard from 'components/RunCard/TransactionRunCard';
import TestHeader from 'components/TestHeader';
import {useTransaction} from 'providers/Transaction/Transaction.provider';
import {useGetTransactionRunsQuery} from 'redux/apis/TraceTest.api';
import {TTransactionRun} from 'types/TransactionRun.types';
import ExperimentalFeature from 'utils/ExperimentalFeature';
import useTransactionCrud from 'providers/Transaction/hooks/useTransactionCrud';
import * as S from './Transaction.styled';

const Content = () => {
  const navigate = useNavigate();
  const {onDelete, transaction} = useTransaction();
  const {runTransaction, isEditLoading} = useTransactionCrud();
  const params = useMemo(() => ({transactionId: transaction.id}), [transaction.id]);

  const handleRunTest = useCallback(async () => {
    if (transaction.id) runTransaction(transaction.id);
  }, [runTransaction, transaction.id]);

  return (
    <S.Container $isWhite={!ExperimentalFeature.isEnabled('transactions')}>
      <TestHeader
        description={transaction.description}
        id={transaction.id}
        onBack={() => navigate('/')}
        onDelete={() => onDelete(transaction.id, transaction.name)}
        title={`${transaction.name} (v${transaction.version})`}
      />

      <S.ActionsContainer>
        <div />
        <Button onClick={handleRunTest} loading={isEditLoading} type="primary" ghost>
          Run Transaction
        </Button>
      </S.ActionsContainer>

      <PaginatedList<TTransactionRun, {transactionId: string}>
        itemComponent={({item}) => (
          <TransactionRunCard
            linkTo={`/transaction/${transaction.id}/run/${item.id}`}
            run={item}
            transactionId={transaction.id}
          />
        )}
        params={params}
        query={useGetTransactionRunsQuery}
      />
    </S.Container>
  );
};

export default Content;
