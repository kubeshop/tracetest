import CreateTestModal from 'components/CreateTestModal/CreateTestModal';
import Empty from 'components/Empty';
import Pagination from 'components/Pagination';
import TestCard from 'components/ResourceCard/TestCard';
import {SortBy, SortDirection, sortOptions} from 'constants/Test.constants';
import usePagination from 'hooks/usePagination';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {useCallback, useState} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import {ADD_TEST_URL, OPENING_TRACETEST_URL} from 'constants/Common.constants';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import useDeleteResource from 'hooks/useDeleteResource';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import VariableSetSelector from 'components/VariableSetSelector/VariableSetSelector';
import Test from 'models/Test.model';
import * as S from './Home.styled';
import CreateButton from './CreateButton';
import HomeFilters from './HomeFilters';
import Loading from './Loading';

const {useGetTestListQuery} = TracetestAPI.instance;

const {onTestClick} = HomeAnalyticsService;
type TParameters = {sortBy: SortBy; sortDirection: SortDirection};
const [{params: defaultSort}] = sortOptions;

const Tests = () => {
  const [isCreateTestOpen, setIsCreateTestOpen] = useState(false);
  const [parameters, setParameters] = useState<TParameters>(defaultSort);

  const pagination = usePagination<Test, TParameters>(useGetTestListQuery, {
    ...parameters,
  });
  const onDelete = useDeleteResource();
  const {runTest} = useTestCrud();
  const {navigate} = useDashboard();

  const handleOnRun = useCallback(
    (test: Test) => {
      runTest({test});
    },
    [runTest]
  );

  const handleOnViewAll = useCallback((id: string) => {
    onTestClick(id);
  }, []);

  const handleOnEdit = useCallback(
    (id: string, lastRunId: number) => {
      navigate(`/test/${id}/run/${lastRunId}`);
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
          <CreateButton onCreate={() => setIsCreateTestOpen(true)} />
        </S.ActionsContainer>

        <Pagination<Test>
          emptyComponent={
            <Empty
              title="Haven't Created a Test Yet"
              message={
                <>
                  Hit the &apos;Create&apos; button below to kickstart your testing adventure. Want to learn more about
                  tests? Just click{' '}
                  <S.Link href={ADD_TEST_URL} target="_blank">
                    here
                  </S.Link>
                  . If you don’t have an app that’s generating OpenTelemetry traces we have a demo for you. Follow these{' '}
                  <S.Link href={OPENING_TRACETEST_URL} target="_blank">
                    instructions
                  </S.Link>
                  !
                </>
              }
              action={<CreateButton onCreate={() => setIsCreateTestOpen(true)} title="Create Your First Test" />}
            />
          }
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
                test={test}
              />
            ))}
          </S.TestListContainer>
        </Pagination>
      </S.Wrapper>

      <CreateTestModal isOpen={isCreateTestOpen} onClose={() => setIsCreateTestOpen(false)} />
    </>
  );
};

export default Tests;
