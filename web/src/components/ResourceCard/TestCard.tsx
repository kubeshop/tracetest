import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {useMemo} from 'react';
import CreateButton from 'components/CreateButton';
import TestRunCard from 'components/RunCard/TestRunCard';
import TracetestAPI from 'redux/apis/Tracetest';
import {ResourceType} from 'types/Resource.type';
import Test from 'models/Test.model';
import TestRun from 'models/TestRun.model';
import * as S from './ResourceCard.styled';
import ResourceCardActions from './ResourceCardActions';
import ResourceCardRuns from './ResourceCardRuns';
import ResourceCardSummary from './ResourceCardSummary';
import useRuns from './useRuns';

const {useLazyGetRunListQuery} = TracetestAPI.instance;

interface IProps {
  onEdit(id: string, lastRunId: number, type: ResourceType): void;
  onDelete(id: string, name: string, type: ResourceType): void;
  onRun(test: Test, type: ResourceType): void;
  onDuplicate(test: Test, type: ResourceType): void;
  onViewAll(id: string, type: ResourceType): void;
  test: Test;
}

const TestCard = ({onEdit, onDelete, onDuplicate, onRun, onViewAll, test}: IProps) => {
  const queryParams = useMemo(() => ({take: 5, testId: test.id}), [test.id]);
  const {isCollapsed, isLoading, list, onClick} = useRuns<TestRun, {testId: string}>(
    useLazyGetRunListQuery,
    queryParams
  );

  const shouldEdit = test.summary.hasRuns;
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

        <S.Row $gap={12}>
          {Test.shouldAllowRun(test.trigger.type) && (
            <CreateButton
              data-cy={`test-run-button-${test.id}`}
              ghost
              onClick={event => {
                event.stopPropagation();
                onRun(test, ResourceType.Test);
              }}
              type="primary"
            >
              Run
            </CreateButton>
          )}
          <ResourceCardActions
            id={test.id}
            shouldEdit={shouldEdit}
            onDelete={() => onDelete(test.id, test.name, ResourceType.Test)}
            onEdit={() => onEdit(test.id, lastRunId, ResourceType.Test)}
            onDuplicate={() => onDuplicate(test, ResourceType.Test)}
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
