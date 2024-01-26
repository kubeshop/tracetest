import {useCallback, useMemo} from 'react';
import CreateButton from 'components/CreateButton';
import PaginatedList from 'components/PaginatedList';
import TestRunCard from 'components/RunCard/TestRunCard';
import TestHeader from 'components/TestHeader';
import useDeleteResource from 'hooks/useDeleteResource';
import useDocumentTitle from 'hooks/useDocumentTitle';
import Test from 'models/Test.model';
import TestRun from 'models/TestRun.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useTest} from 'providers/Test/Test.provider';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import TracetestAPI from 'redux/apis/Tracetest';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {ResourceType} from 'types/Resource.type';
import * as S from './Test.styled';

const {useGetRunListQuery} = TracetestAPI.instance;

const Content = () => {
  const {test} = useTest();
  const onDeleteResource = useDeleteResource();
  const {runTest, isLoadingRunTest, duplicate} = useTestCrud();
  const {onOpen} = useConfirmationModal();
  const params = useMemo(() => ({testId: test.id}), [test.id]);
  useDocumentTitle(`${test.name}`);

  const {navigate} = useDashboard();

  const shouldEdit = test.summary.hasRuns;
  const onEdit = () => navigate(`/test/${test.id}/run/${test.summary.runs}`);

  const handleOnDuplicate = useCallback(() => {
    onOpen({
      heading: `Duplicate Test`,
      title: `Create a duplicated version of Test: ${test.name}`,
      okText: 'Duplicate',
      onConfirm: () => duplicate(test),
    });
  }, [duplicate, onOpen, test]);

  return (
    <S.Container $isWhite>
      <TestHeader
        description={`${test.trigger.type.toUpperCase()} • ${test.trigger.method.toUpperCase()} • ${
          test.trigger.entryPoint
        }`}
        id={test.id}
        onBack={() => navigate('/tests')}
        onDelete={() => onDeleteResource(test.id, test.name, ResourceType.Test)}
        onEdit={onEdit}
        onDuplicate={handleOnDuplicate}
        shouldEdit={shouldEdit}
        title={`${test.name} (v${test.version})`}
        runButton={
          Test.shouldAllowRun(test.trigger.type) ? (
            <CreateButton
              data-cy="test-details-run-test-button"
              ghost
              loading={isLoadingRunTest}
              onClick={() => runTest({test})}
              type="primary"
            >
              Run Test
            </CreateButton>
          ) : null
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
