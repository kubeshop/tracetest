import {useCallback, useMemo} from 'react';
import CreateButton from 'components/CreateButton';
import PaginatedList from 'components/PaginatedList';
import TestSuiteRunCard from 'components/RunCard/TestSuiteRunCard';
import TestHeader from 'components/TestHeader';
import useDocumentTitle from 'hooks/useDocumentTitle';
import TestSuiteRun from 'models/TestSuiteRun.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useTestSuite} from 'providers/TestSuite/TestSuite.provider';
import {useTestSuiteCrud} from 'providers/TestSuite';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import TracetestAPI from 'redux/apis/Tracetest';
import * as S from './TestSuite.styled';

const {useGetTestSuiteRunsQuery} = TracetestAPI.instance;

const Content = () => {
  const {onDelete, testSuite} = useTestSuite();
  const {runTestSuite, duplicate, isEditLoading} = useTestSuiteCrud();
  const {onOpen} = useConfirmationModal();
  const params = useMemo(() => ({testSuiteId: testSuite.id}), [testSuite.id]);

  useDocumentTitle(`${testSuite.name}`);

  const handleRunTest = useCallback(async () => {
    if (testSuite.id) runTestSuite(testSuite);
  }, [runTestSuite, testSuite]);

  const {navigate} = useDashboard();

  const shouldEdit = testSuite.summary.hasRuns;
  const onEdit = () => navigate(`/testsuite/${testSuite.id}/run/${testSuite.summary.runs}`);

  const handleOnDuplicate = useCallback(() => {
    onOpen({
      heading: `Duplicate Test Suite`,
      title: `Create a duplicated version of Test Suite: ${testSuite.name}`,
      okText: 'Duplicate',
      onConfirm: () => duplicate(testSuite),
    });
  }, [duplicate, onOpen, testSuite]);

  return (
    <S.Container $isWhite>
      <TestHeader
        description={testSuite.description}
        id={testSuite.id}
        onDelete={() => onDelete(testSuite.id, testSuite.name)}
        onEdit={onEdit}
        onDuplicate={handleOnDuplicate}
        shouldEdit={shouldEdit}
        title={`${testSuite.name} (v${testSuite.version})`}
        runButton={
          <CreateButton ghost loading={isEditLoading} onClick={handleRunTest} type="primary">
            Run Test Suite
          </CreateButton>
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
