import AllowButton, {Operation} from 'components/AllowButton';
import CreateButton from 'components/CreateButton';
import Pagination from 'components/Pagination';
import TestCard from 'components/ResourceCard/TestCard';
import TriggerTypeModal from 'components/TriggerTypeModal';
import {SortBy, SortDirection, sortOptions} from 'constants/Test.constants';
import usePagination from 'hooks/usePagination';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {useCallback, useState} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import useDeleteResource from 'hooks/useDeleteResource';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import VariableSetSelector from 'components/VariableSetSelector/VariableSetSelector';
import Test from 'models/Test.model';
import ImportModal from 'components/ImportModal/ImportModal';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import * as S from './Home.styled';
import HomeFilters from './HomeFilters';
import Loading from './Loading';
import EmptyTestList from './EmptyTestList';

const {useGetTestListQuery} = TracetestAPI.instance;

const {onTestClick} = HomeAnalyticsService;
type TParameters = {sortBy: SortBy; sortDirection: SortDirection};
const [{params: defaultSort}] = sortOptions;

const Tests = () => {
  const [isCreateTestOpen, setIsCreateTestOpen] = useState(false);
  const [isImportTestOpen, setIsImportTestOpen] = useState(false);
  const [parameters, setParameters] = useState<TParameters>(defaultSort);
  const {onOpen} = useConfirmationModal();

  const pagination = usePagination<Test, TParameters>(useGetTestListQuery, {
    ...parameters,
  });
  const onDelete = useDeleteResource();
  const {runTest, duplicate} = useTestCrud();
  const {navigate} = useDashboard();

  const handleOnRun = useCallback(
    (test: Test) => {
      runTest({test});
    },
    [runTest]
  );

  const handleOnDuplicate = useCallback(
    (test: Test) => {
      onOpen({
        heading: `Duplicate Test`,
        title: `Create a duplicated version of Test: ${test.name}`,
        okText: 'Duplicate',
        onConfirm: () => duplicate(test),
      });
    },
    [duplicate, onOpen]
  );

  const handleOnViewAll = useCallback((id: string) => {
    onTestClick(id);
  }, []);

  const handleOnEdit = useCallback(
    (id: string, lastRunId: number) => {
      navigate(`/test/${id}/run/${lastRunId}/trigger`);
    },
    [navigate]
  );

  return (
    <>
      <S.Wrapper>
        <S.HeaderContainer>
          <S.TitleText>All Tests</S.TitleText>
          <VariableSetSelector />
        </S.HeaderContainer>

        <S.ActionsContainer>
          <HomeFilters
            onSearch={pagination.search}
            onSortBy={(sortBy, sortDirection) => setParameters({sortBy, sortDirection})}
            isEmpty={pagination.list?.length === 0}
          />
          <S.ButtonsContainer>
            <AllowButton
              operation={Operation.Edit}
              ButtonComponent={CreateButton}
              data-cy="import-button"
              onClick={() => setIsImportTestOpen(true)}
              type="primary"
              ghost
            >
              Import
            </AllowButton>

            <AllowButton
              operation={Operation.Edit}
              ButtonComponent={CreateButton}
              data-cy="create-button"
              onClick={() => setIsCreateTestOpen(true)}
              type="primary"
            >
              Create
            </AllowButton>
          </S.ButtonsContainer>
        </S.ActionsContainer>

        <Pagination<Test>
          emptyComponent={<EmptyTestList onClick={() => setIsCreateTestOpen(true)} />}
          loadingComponent={<Loading />}
          {...pagination}
        >
          <S.TestListContainer data-cy="test-list">
            {pagination.list?.map(test => (
              <TestCard
                key={test.id}
                onEdit={handleOnEdit}
                onDelete={onDelete}
                onRun={handleOnRun}
                onViewAll={handleOnViewAll}
                onDuplicate={handleOnDuplicate}
                test={test}
              />
            ))}
          </S.TestListContainer>
        </Pagination>
      </S.Wrapper>

      <TriggerTypeModal isOpen={isCreateTestOpen} onClose={() => setIsCreateTestOpen(false)} />
      <ImportModal isOpen={isImportTestOpen} onClose={() => setIsImportTestOpen(false)} />
    </>
  );
};

export default Tests;
