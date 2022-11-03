import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {Skeleton} from 'antd';
import {useCallback, useState} from 'react';

import {useLazyGetRunListQuery} from 'redux/apis/TraceTest.api';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {ResourceType} from 'types/Resource.type';
import {TTest} from 'types/Test.types';
import ResultCardList from '../RunCardList/RunCardList';
import * as S from './ResourceCard.styled';
import ResourceCardActions from './ResourceCardActions';
import ResourceCardSummary from './ResourceCardSummary';

interface IProps {
  onDelete(id: string, name: string, type: ResourceType): void;
  onRun(id: string, type: ResourceType): void;
  onViewAll(id: string, type: ResourceType): void;
  test: TTest;
}

const TestCard = ({onDelete, onRun, onViewAll, test}: IProps) => {
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
    getRuns({testId: test.id, take: 5});
  }, [getRuns, isCollapsed, runs.length, test.id]);

  return (
    <S.Container>
      <S.TestContainer onClick={() => handleOnClick()}>
        {isCollapsed ? <DownOutlined /> : <RightOutlined data-cy={`collapse-test-${test.id}`} />}
        <S.Box $type={ResourceType.test}>
          <S.BoxTitle level={2}>{test.summary.runs}</S.BoxTitle>
        </S.Box>
        <S.TitleContainer>
          <S.Title level={3}>{test.name}</S.Title>
          <S.Text>
            {test.trigger.method} • {test.trigger.entryPoint}
          </S.Text>
        </S.TitleContainer>

        <ResourceCardSummary summary={test.summary} />

        <S.Row>
          <S.RunButton
            type="primary"
            ghost
            data-cy={`test-run-button-${test.id}`}
            onClick={event => {
              event.stopPropagation();
              onRun(test.id, ResourceType.test);
            }}
          >
            Run
          </S.RunButton>
          <ResourceCardActions id={test.id} onDelete={() => onDelete(test.id, test.name, ResourceType.test)} />
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

          {Boolean(runs.length) && <ResultCardList testId={test.id} resultList={runs} />}

          {runs.length === 5 && (
            <S.FooterContainer>
              <S.Link data-cy="test-details-link" onClick={() => onViewAll(test.id, ResourceType.test)}>
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

export default TestCard;
