import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {useMemo} from 'react';

import {useLazyGetRunListQuery} from 'redux/apis/TraceTest.api';
import {ResourceType} from 'types/Resource.type';
import {TTestRun} from 'types/TestRun.types';
import {TTransaction} from 'types/Transaction.types';
import ResultCardList from '../RunCardList/RunCardList';
import * as S from './ResourceCard.styled';
import ResourceCardActions from './ResourceCardActions';
import ResourceCardRuns from './ResourceCardRuns';
import ResourceCardSummary from './ResourceCardSummary';
import useRuns from './useRuns';

interface IProps {
  onDelete(id: string, name: string, type: ResourceType): void;
  onRun(id: string, type: ResourceType): void;
  onViewAll(id: string, type: ResourceType): void;
  transaction: TTransaction;
}

const TransactionCard = ({onDelete, onRun, onViewAll, transaction}: IProps) => {
  const queryParams = useMemo(() => ({take: 5, testId: transaction.id}), [transaction.id]);
  const {isCollapsed, isLoading, list, onClick} = useRuns<TTestRun, {testId: string}>(
    useLazyGetRunListQuery,
    queryParams
  );

  return (
    <S.Container $type={ResourceType.Transaction}>
      <S.TestContainer onClick={onClick}>
        {isCollapsed ? <RightOutlined data-cy={`collapse-transaction-${transaction.id}`} /> : <DownOutlined />}
        <S.Box $type={ResourceType.Transaction}>
          <S.BoxTitle level={2}>1</S.BoxTitle>
        </S.Box>
        <S.TitleContainer>
          <S.Title level={3}>{transaction.name}</S.Title>
          <S.Text>{transaction.description}</S.Text>
        </S.TitleContainer>

        <ResourceCardSummary />

        <S.Row>
          <S.RunButton
            type="primary"
            ghost
            data-cy={`transaction-run-button-${transaction.id}`}
            onClick={event => {
              event.stopPropagation();
              onRun(transaction.id, ResourceType.Transaction);
            }}
          >
            Run
          </S.RunButton>
          <ResourceCardActions
            id={transaction.id}
            onDelete={() => onDelete(transaction.id, transaction.name, ResourceType.Transaction)}
          />
        </S.Row>
      </S.TestContainer>

      <ResourceCardRuns
        hasMoreRuns={list.length === 5}
        hasRuns={Boolean(list.length)}
        isCollapsed={isCollapsed}
        isLoading={isLoading}
        onViewAll={() => onViewAll(transaction.id, ResourceType.Test)}
      >
        <ResultCardList testId={transaction.id} resultList={list} />
      </ResourceCardRuns>
    </S.Container>
  );
};

export default TransactionCard;
