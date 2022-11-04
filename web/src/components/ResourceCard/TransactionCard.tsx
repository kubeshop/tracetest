import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {Skeleton} from 'antd';
import {useCallback, useState} from 'react';

import {useLazyGetRunListQuery} from 'redux/apis/TraceTest.api';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {TTransaction} from 'types/Transaction.types';
import {ResourceType} from 'types/Resource.type';
import ResultCardList from '../RunCardList/RunCardList';
import * as S from './ResourceCard.styled';
import ResourceCardActions from './ResourceCardActions';
import ResourceCardSummary from './ResourceCardSummary';

interface IProps {
  onDelete(id: string, name: string, type: ResourceType): void;
  onRun(id: string, type: ResourceType): void;
  onViewAll(id: string, type: ResourceType): void;
  transaction: TTransaction;
}

const TransactionCard = ({onDelete, onRun, onViewAll, transaction}: IProps) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [getRuns, {data, isLoading}] = useLazyGetRunListQuery();
  const runs = data?.items || [];

  const handleOnClick = useCallback(() => {
    if (isCollapsed) {
      setIsCollapsed(false);
      return;
    }

    setIsCollapsed(true);
    TestAnalyticsService.onTestCardCollapse();
    if (runs.length > 0) {
      return;
    }
    getRuns({testId: transaction.id, take: 5});
  }, [getRuns, isCollapsed, runs.length, transaction.id]);

  return (
    <S.Container $type={ResourceType.transaction}>
      <S.TestContainer onClick={() => handleOnClick()}>
        {isCollapsed ? <DownOutlined /> : <RightOutlined data-cy={`collapse-transaction-${transaction.id}`} />}
        <S.Box $type={ResourceType.transaction}>
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
              onRun(transaction.id, ResourceType.transaction);
            }}
          >
            Run
          </S.RunButton>
          <ResourceCardActions
            id={transaction.id}
            onDelete={() => onDelete(transaction.id, transaction.name, ResourceType.transaction)}
          />
        </S.Row>
      </S.TestContainer>

      {isCollapsed && (
        <S.RunsContainer>
          {isLoading && (
            <S.LoadingContainer direction="vertical">
              <Skeleton.Input active block size="small" />
              <Skeleton.Input active block size="small" />
              <Skeleton.Input active block size="small" />
            </S.LoadingContainer>
          )}

          {Boolean(runs.length) && <ResultCardList testId={transaction.id} resultList={runs} />}

          {runs.length === 5 && (
            <S.FooterContainer>
              <S.Link data-cy="test-details-link" onClick={() => onViewAll(transaction.id, ResourceType.transaction)}>
                View all runs
              </S.Link>
            </S.FooterContainer>
          )}

          {!runs.length && !isLoading && (
            <S.EmptyStateContainer>
              <S.EmptyStateIcon />
              <S.Text disabled>No Runs</S.Text>
            </S.EmptyStateContainer>
          )}
        </S.RunsContainer>
      )}
    </S.Container>
  );
};

export default TransactionCard;
