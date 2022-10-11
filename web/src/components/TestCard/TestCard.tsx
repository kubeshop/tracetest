import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {Skeleton, Tooltip} from 'antd';
import {useCallback, useState} from 'react';

import {useLazyGetRunListQuery} from 'redux/apis/TraceTest.api';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {TTest} from 'types/Test.types';
import Date from 'utils/Date';
import * as S from './TestCard.styled';
import TestCardActions from './TestCardActions';
import ResultCardList from '../RunCardList/RunCardList';

interface IProps {
  onDelete(test: TTest): void;
  onRun(id: string): void;
  onViewAll(id: string): void;
  test: TTest;
}

const TestCard = ({onDelete, onRun, onViewAll, test}: IProps) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [getRuns, {data: runs = [], isLoading}] = useLazyGetRunListQuery();

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
        <S.Box>
          <S.BoxTitle level={2}>{test.summary.runs}</S.BoxTitle>
        </S.Box>
        <S.TitleContainer>
          <S.Title level={3}>{test.name}</S.Title>
          <S.Text>
            {test.trigger.method} • {test.trigger.entryPoint}
          </S.Text>
        </S.TitleContainer>
        <div>
          <S.Text>Last run time:</S.Text>
          <Tooltip title={Date.format(test.summary.lastRun.time)}>
            <S.Text>{Date.getTimeAgo(test.summary.lastRun.time)}</S.Text>
          </Tooltip>
        </div>

        <div>
          <S.Text>Last run result:</S.Text>
          <S.Row>
            <Tooltip title="Passed assertions">
              <S.HeaderDetail>
                <S.HeaderDot $passed />
                {test.summary.lastRun.passes}
              </S.HeaderDetail>
            </Tooltip>
            <Tooltip title="Failed assertions">
              <S.HeaderDetail>
                <S.HeaderDot $passed={false} />
                {test.summary.lastRun.fails}
              </S.HeaderDetail>
            </Tooltip>
          </S.Row>
        </div>

        <S.Row>
          <S.RunButton
            type="primary"
            ghost
            data-cy={`test-run-button-${test.id}`}
            onClick={event => {
              event.stopPropagation();
              onRun(test.id);
            }}
          >
            Run
          </S.RunButton>
          <TestCardActions testId={test.id} onDelete={() => onDelete(test)} />
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
              <S.Link data-cy="test-details-link" onClick={() => onViewAll(test.id)}>
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
