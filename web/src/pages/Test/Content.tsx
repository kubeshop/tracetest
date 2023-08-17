import {useMemo} from 'react';
import {Button} from 'antd';
import PaginatedList from 'components/PaginatedList';
import TestRunCard from 'components/RunCard/TestRunCard';
import TestHeader from 'components/TestHeader';
import useDeleteResource from 'hooks/useDeleteResource';
import useDocumentTitle from 'hooks/useDocumentTitle';
import TestRun from 'models/TestRun.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useTest} from 'providers/Test/Test.provider';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import TracetestAPI from 'redux/apis/Tracetest';
import {ResourceType} from 'types/Resource.type';
import * as S from './Test.styled';

const {useGetRunListQuery} = TracetestAPI.instance;

const Content = () => {
  const {test} = useTest();
  const onDeleteResource = useDeleteResource();
  const {runTest, isLoadingRunTest} = useTestCrud();
  const params = useMemo(() => ({testId: test.id}), [test.id]);
  useDocumentTitle(`${test.name}`);

  const {navigate} = useDashboard();

  const shouldEdit = test.summary.hasRuns;
  const onEdit = () => navigate(`/test/${test.id}/run/${test.summary.runs}`);

  return (
    <S.Container $isWhite>
      <TestHeader
        description={`${test.trigger.type.toUpperCase()} • ${test.trigger.method.toUpperCase()} • ${
          test.trigger.entryPoint
        }`}
        id={test.id}
        onDelete={() => onDeleteResource(test.id, test.name, ResourceType.Test)}
        onEdit={onEdit}
        shouldEdit={shouldEdit}
        title={`${test.name} (v${test.version})`}
        runButton={
          <Button
            data-cy="test-details-run-test-button"
            ghost
            loading={isLoadingRunTest}
            onClick={() => runTest({test})}
            type="primary"
          >
            Run Test
          </Button>
        }
      />

      <PaginatedList<TestRun, {testId: string}>
        dataCy="run-card-list"
        itemComponent={({item}) => (
          <TestRunCard linkTo={`/test/${test.id}/run/${item.id}`} run={item} testId={test.id} />
        )}
        params={params}
        query={useGetRunListQuery}
      />
    </S.Container>
  );
};

export default Content;
