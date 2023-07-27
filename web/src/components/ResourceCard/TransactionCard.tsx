import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {useMemo} from 'react';

import TransactionRunCard from 'components/RunCard/TransactionRunCard';
import {useLazyGetTransactionRunsQuery} from 'redux/apis/Tracetest';
import {ResourceType} from 'types/Resource.type';
import Transaction from 'models/Transaction.model';
import TransactionRun from 'models/TransactionRun.model';
import * as S from './ResourceCard.styled';
import ResourceCardActions from './ResourceCardActions';
import ResourceCardRuns from './ResourceCardRuns';
import ResourceCardSummary from './ResourceCardSummary';
import useRuns from './useRuns';

interface IProps {
  onEdit(id: string, lastRunId: number, type: ResourceType): void;
  onDelete(id: string, name: string, type: ResourceType): void;
  onRun(transaction: Transaction, type: ResourceType): void;
  onViewAll(id: string, type: ResourceType): void;
  transaction: Transaction;
}

const TransactionCard = ({
  onEdit,
  onDelete,
  onRun,
  onViewAll,
  transaction: {id: transactionId, summary, name, description},
  transaction,
}: IProps) => {
  const queryParams = useMemo(() => ({take: 5, transactionId}), [transactionId]);
  const {isCollapsed, isLoading, list, onClick} = useRuns<TransactionRun, {transactionId: string}>(
    useLazyGetTransactionRunsQuery,
    queryParams
  );

  const shouldEdit = summary.hasRuns;
  const lastRunId = summary.runs; // assume the total of runs as the last run

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

        <ResourceCardSummary summary={summary} />

        <S.Row>
          <S.RunButton
            type="primary"
            ghost
            data-cy={`transaction-run-button-${transactionId}`}
            onClick={event => {
              event.stopPropagation();
              onRun(transaction, ResourceType.Transaction);
            }}
          >
            Run
          </S.RunButton>
          <ResourceCardActions
            id={transactionId}
            shouldEdit={shouldEdit}
            onDelete={() => onDelete(transactionId, name, ResourceType.Transaction)}
            onEdit={() => onEdit(transactionId, lastRunId, ResourceType.Transaction)}
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
