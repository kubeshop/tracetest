import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {useMemo} from 'react';

import CreateButton from 'components/CreateButton';
import TestSuiteRunCard from 'components/RunCard/TestSuiteRunCard';
import TracetestAPI from 'redux/apis/Tracetest';
import {ResourceType} from 'types/Resource.type';
import TestSuite from 'models/TestSuite.model';
import TestSuiteRun from 'models/TestSuiteRun.model';
import Box from './Box';
import * as S from './ResourceCard.styled';
import ResourceCardActions from './ResourceCardActions';
import ResourceCardRuns from './ResourceCardRuns';
import ResourceCardSummary from './ResourceCardSummary';
import useRuns from './useRuns';

const {useLazyGetTestSuiteRunsQuery} = TracetestAPI.instance;

interface IProps {
  onEdit(id: string, lastRunId: number, type: ResourceType): void;
  onDelete(id: string, name: string, type: ResourceType): void;
  onRun(testSuite: TestSuite, type: ResourceType): void;
  onDuplicate(testSuite: TestSuite): void;
  onViewAll(id: string, type: ResourceType): void;
  testSuite: TestSuite;
}

const TestSuiteCard = ({
  onEdit,
  onDelete,
  onRun,
  onDuplicate,
  onViewAll,
  testSuite: {id: testSuiteId, summary, name, description},
  testSuite,
}: IProps) => {
  const queryParams = useMemo(() => ({take: 5, testSuiteId}), [testSuiteId]);
  const {isCollapsed, isLoading, list, onClick} = useRuns<TestSuiteRun, {testSuiteId: string}>(
    useLazyGetTestSuiteRunsQuery,
    queryParams
  );

  const shouldEdit = summary.hasRuns;
  const lastRunId = summary.runs; // assume the total of runs as the last run

  return (
    <S.Container $type={ResourceType.TestSuite}>
      <S.TestContainer onClick={onClick}>
        {isCollapsed ? <RightOutlined data-cy={`collapse-testsuite-${testSuiteId}`} /> : <DownOutlined />}

        <Box num={summary.runs} type={ResourceType.TestSuite} />

        <S.TitleContainer>
          <S.Title level={3}>{name}</S.Title>
          <S.Text>{description}</S.Text>
        </S.TitleContainer>

        <ResourceCardSummary summary={summary} />

        <S.Row $gap={12}>
          <CreateButton
            data-cy={`testsuite-run-button-${testSuiteId}`}
            ghost
            onClick={event => {
              event.stopPropagation();
              onRun(testSuite, ResourceType.TestSuite);
            }}
            type="primary"
          >
            Run
          </CreateButton>
          <ResourceCardActions
            id={testSuiteId}
            shouldEdit={shouldEdit}
            onDelete={() => onDelete(testSuiteId, name, ResourceType.TestSuite)}
            onEdit={() => onEdit(testSuiteId, lastRunId, ResourceType.TestSuite)}
            onDuplicate={() => onDuplicate(testSuite)}
          />
        </S.Row>
      </S.TestContainer>

      <ResourceCardRuns
        hasMoreRuns={list.length === 5}
        hasRuns={Boolean(list.length)}
        isCollapsed={isCollapsed}
        isLoading={isLoading}
        resourcePath={`/testsuite/${testSuiteId}`}
        onViewAll={() => onViewAll(testSuiteId, ResourceType.TestSuite)}
      >
        <S.RunsListContainer>
          {list.map(run => (
            <TestSuiteRunCard
              key={run.id}
              linkTo={`/testsuite/${testSuiteId}/run/${run.id}`}
              run={run}
              testSuiteId={testSuiteId}
            />
          ))}
        </S.RunsListContainer>
      </ResourceCardRuns>
    </S.Container>
  );
};

export default TestSuiteCard;
