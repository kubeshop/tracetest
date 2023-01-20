import {Button} from 'antd';
import {useCallback, useMemo} from 'react';
import {useNavigate} from 'react-router-dom';

import PaginatedList from 'components/PaginatedList';
import TransactionRunCard from 'components/RunCard/TransactionRunCard';
import TestHeader from 'components/TestHeader';
import {useTransaction} from 'providers/Transaction/Transaction.provider';
import {useGetTransactionRunsQuery} from 'redux/apis/TraceTest.api';
import {TTransactionRun} from 'types/TransactionRun.types';
import useTransactionCrud from 'providers/Transaction/hooks/useTransactionCrud';
import useDocumentTitle from 'hooks/useDocumentTitle';
import * as S from './Transaction.styled';

const Content = () => {
  const {onDelete, transaction} = useTransaction();
  const {runTransaction, isEditLoading} = useTransactionCrud();
  const params = useMemo(() => ({transactionId: transaction.id}), [transaction.id]);

  useDocumentTitle(`${transaction.name}`);

  const handleRunTest = useCallback(async () => {
    if (transaction.id) runTransaction(transaction);
  }, [runTransaction, transaction]);

  const navigate = useNavigate();
  const canEdit = test != null;
  const onEdit = () => navigate(`/transaction/${transaction.id}/run/${transaction.version}`);

  return (
    <S.Container $isWhite>
      <TestHeader
        description={transaction.description}
        id={transaction.id}
        onDelete={() => onDelete(transaction.id, transaction.name)}
        onEdit={onEdit}
        canEdit={canEdit}
        title={`${transaction.name} (v${transaction.version})`}
        runButton={
          <Button onClick={handleRunTest} loading={isEditLoading} type="primary" ghost>
            Run Transaction
          </Button>
        }
      />

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
