import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {useMemo} from 'react';

import TransactionRunCard from 'components/RunCard/TransactionRunCard';
import {useLazyGetTransactionRunsQuery} from 'redux/apis/TraceTest.api';
import {ResourceType} from 'types/Resource.type';
import {TTransactionRun} from 'types/TransactionRun.types';
import {TTransaction} from 'types/Transaction.types';
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

const TransactionCard = ({
  onDelete,
  onRun,
  onViewAll,
  transaction: {id: transactionId, summary, name, description},
}: IProps) => {
  const queryParams = useMemo(() => ({take: 5, transactionId}), [transactionId]);
  const {isCollapsed, isLoading, list, onClick} = useRuns<TTransactionRun, {transactionId: string}>(
    useLazyGetTransactionRunsQuery,
    queryParams
  );

  return (
    <S.Container $type={ResourceType.Transaction}>
      <S.TestContainer onClick={onClick}>
        {isCollapsed ? <RightOutlined data-cy={`collapse-transaction-${transactionId}`} /> : <DownOutlined />}
        <S.Box $type={ResourceType.Transaction}>
          <S.BoxTitle level={2}>{summary.runs}</S.BoxTitle>
        </S.Box>
        <S.TitleContainer>
          <S.Title level={3}>{name}</S.Title>
          <S.Text>{description}</S.Text>
        </S.TitleContainer>

        <ResourceCardSummary summary={summary} shouldShowResult={false} />

        <S.Row>
          <S.RunButton
            type="primary"
            ghost
            data-cy={`transaction-run-button-${transactionId}`}
            onClick={event => {
              event.stopPropagation();
              onRun(transactionId, ResourceType.Transaction);
            }}
          >
            Run
          </S.RunButton>
          <ResourceCardActions
            id={transactionId}
            onDelete={() => onDelete(transactionId, name, ResourceType.Transaction)}
          />
        </S.Row>
      </S.TestContainer>

      <ResourceCardRuns
        hasMoreRuns={list.length === 5}
        hasRuns={Boolean(list.length)}
        isCollapsed={isCollapsed}
        isLoading={isLoading}
        resourcePath={`/transaction/${transactionId}`}
        onViewAll={() => onViewAll(transactionId, ResourceType.Transaction)}
      >
        <S.RunsListContainer>
          {list.map(run => (
            <TransactionRunCard
              key={run.id}
              linkTo={`/transaction/${transactionId}/run/${run.id}`}
              run={run}
              transactionId={transactionId}
            />
          ))}
        </S.RunsListContainer>
      </ResourceCardRuns>
    </S.Container>
  );
};

export default TransactionCard;
