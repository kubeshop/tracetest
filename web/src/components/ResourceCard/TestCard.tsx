import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {useMemo} from 'react';

import TestRunCard from 'components/RunCard/TestRunCard';
import {useLazyGetRunListQuery} from 'redux/apis/TraceTest.api';
import {ResourceType} from 'types/Resource.type';
import {TTest} from 'types/Test.types';
import {TTestRun} from 'types/TestRun.types';
import * as S from './ResourceCard.styled';
import ResourceCardActions from './ResourceCardActions';
import ResourceCardRuns from './ResourceCardRuns';
import ResourceCardSummary from './ResourceCardSummary';
import useRuns from './useRuns';

interface IProps {
  onEdit(id: string, lastRunId: number, type: ResourceType): void;
  onDelete(id: string, name: string, type: ResourceType): void;
  onRun(test: TTest, type: ResourceType): void;
  onViewAll(id: string, type: ResourceType): void;
  test: TTest;
}

const TestCard = ({onEdit, onDelete, onRun, onViewAll, test}: IProps) => {
  const queryParams = useMemo(() => ({take: 5, testId: test.id}), [test.id]);
  const {isCollapsed, isLoading, list, onClick} = useRuns<TTestRun, {testId: string}>(
    useLazyGetRunListQuery,
    queryParams
  );

  const canEdit = test.summary.runs > 0;
  const lastRunId = test.summary.runs; // assume the total of runs as the last run

  return (
    <S.Container $type={ResourceType.Test}>
      <S.TestContainer onClick={onClick}>
        {isCollapsed ? <RightOutlined data-cy={`collapse-test-${test.id}`} /> : <DownOutlined />}
        <S.Box $type={ResourceType.Test}>
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
              onRun(test, ResourceType.Test);
            }}
          >
            Run
          </S.RunButton>
          <ResourceCardActions
            id={test.id}
            canEdit={canEdit}
            onDelete={() => onDelete(test.id, test.name, ResourceType.Test)}
            onEdit={() => onEdit(test.id, lastRunId, ResourceType.Test)}
           />
        </S.Row>
      </S.TestContainer>

      <ResourceCardRuns
        hasMoreRuns={list.length === 5}
        hasRuns={Boolean(list.length)}
        isCollapsed={isCollapsed}
        isLoading={isLoading}
        onViewAll={() => onViewAll(test.id, ResourceType.Test)}
        resourcePath={`/test/${test.id}`}
      >
        <S.RunsListContainer data-cy="run-card-list">
          {list.map(run => (
            <TestRunCard key={run.id} linkTo={`/test/${test.id}/run/${run.id}`} run={run} testId={test.id} />
          ))}
        </S.RunsListContainer>
      </ResourceCardRuns>
    </S.Container>
  );
};

export default TestCard;
