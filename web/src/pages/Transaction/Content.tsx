import {Button} from 'antd';
import {useMemo} from 'react';
import {useNavigate} from 'react-router-dom';

import PaginatedList from 'components/PaginatedList';
import RunCard from 'components/RunCard';
import TestHeader from 'components/TestHeader';
import {useTransaction} from 'providers/Transaction/Transaction.provider';
import {useGetRunListQuery} from 'redux/apis/TraceTest.api';
import {TTestRun} from 'types/TestRun.types';
import ExperimentalFeature from 'utils/ExperimentalFeature';
import * as S from './Transaction.styled';

const Content = () => {
  const navigate = useNavigate();
  const {isLoadingRun, onDelete, onRun, transaction} = useTransaction();
  const params = useMemo(() => ({testId: transaction.id}), [transaction.id]);

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
        <Button onClick={onRun} loading={isLoadingRun} type="primary" ghost>
          Run Transaction
        </Button>
      </S.ActionsContainer>

      <PaginatedList<TTestRun, {testId: string}>
        itemComponent={({item}) => (
          <RunCard linkTo={`/test/${transaction.id}/run/${item.id}`} run={item} testId="123" />
        )}
        params={params}
        query={useGetRunListQuery}
      />
    </S.Container>
  );
};

export default Content;
