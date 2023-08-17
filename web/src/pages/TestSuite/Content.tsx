import {Button} from 'antd';
import {useCallback, useMemo} from 'react';
import PaginatedList from 'components/PaginatedList';
import TestSuiteRunCard from 'components/RunCard/TestSuiteRunCard';
import TestHeader from 'components/TestHeader';
import useDocumentTitle from 'hooks/useDocumentTitle';
import TestSuiteRun from 'models/TestSuiteRun.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useTestSuite} from 'providers/TestSuite/TestSuite.provider';
import {useTestSuiteCrud} from 'providers/TestSuite';
import TracetestAPI from 'redux/apis/Tracetest';
import * as S from './TestSuite.styled';

const {useGetTestSuiteRunsQuery} = TracetestAPI.instance;

const Content = () => {
  const {onDelete, testSuite} = useTestSuite();
  const {runTestSuite, isEditLoading} = useTestSuiteCrud();
  const params = useMemo(() => ({testSuiteId: testSuite.id}), [testSuite.id]);

  useDocumentTitle(`${testSuite.name}`);

  const handleRunTest = useCallback(async () => {
    if (testSuite.id) runTestSuite(testSuite);
  }, [runTestSuite, testSuite]);

  const {navigate} = useDashboard();

  const shouldEdit = testSuite.summary.hasRuns;
  const onEdit = () => navigate(`/testsuite/${testSuite.id}/run/${testSuite.summary.runs}`);

  return (
    <S.Container $isWhite>
      <TestHeader
        description={testSuite.description}
        id={testSuite.id}
        onDelete={() => onDelete(testSuite.id, testSuite.name)}
        onEdit={onEdit}
        shouldEdit={shouldEdit}
        title={`${testSuite.name} (v${testSuite.version})`}
        runButton={
          <Button onClick={handleRunTest} loading={isEditLoading} type="primary" ghost>
            Run Test Suite
          </Button>
        }
      />

      <PaginatedList<TestSuiteRun, {testSuiteId: string}>
        itemComponent={({item}) => (
          <TestSuiteRunCard
            linkTo={`/testsuite/${testSuite.id}/run/${item.id}`}
            run={item}
            testSuiteId={testSuite.id}
          />
        )}
        params={params}
        query={useGetTestSuiteRunsQuery}
      />
    </S.Container>
  );
};

export default Content;
