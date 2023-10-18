import AllowButton, {Operation} from 'components/AllowButton';
import CreateButton from 'components/CreateButton';
import CreateTestSuiteModal from 'components/CreateTestSuiteModal/CreateTestSuiteModal';
import Empty from 'components/Empty';
import Pagination from 'components/Pagination';
import TestSuiteCard from 'components/ResourceCard/TestSuiteCard';
import {SortBy, SortDirection, sortOptions} from 'constants/Test.constants';
import useDeleteResource from 'hooks/useDeleteResource';
import usePagination from 'hooks/usePagination';
import {useCallback, useState} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import {ADD_TEST_SUITE_URL} from 'constants/Common.constants';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import useTestSuiteCrud from 'providers/TestSuite/hooks/useTestSuiteCrud';
import VariableSetSelector from 'components/VariableSetSelector/VariableSetSelector';
import TestSuite from 'models/TestSuite.model';
import * as S from './TestSuites.styled';
import HomeFilters from '../Home/HomeFilters';
import Loading from '../Home/Loading';

const {useGetTestListQuery, useGetTestSuiteListQuery} = TracetestAPI.instance;

const {onTestClick} = HomeAnalyticsService;
type TParameters = {sortBy: SortBy; sortDirection: SortDirection};
const [{params: defaultSort}] = sortOptions;

const Resources = () => {
  const [isCreateTestSuiteOpen, setIsCreateTestSuiteOpen] = useState(false);
  const [parameters, setParameters] = useState<TParameters>(defaultSort);

  const {data: testListData} = useGetTestListQuery({});
  const pagination = usePagination<TestSuite, TParameters>(useGetTestSuiteListQuery, parameters);
  const onDelete = useDeleteResource();
  const {runTestSuite} = useTestSuiteCrud();
  const {navigate} = useDashboard();

  const handleOnRun = useCallback(
    (resource: TestSuite) => {
      runTestSuite(resource as TestSuite);
    },
    [runTestSuite]
  );

  const handleOnViewAll = useCallback((id: string) => {
    onTestClick(id);
  }, []);

  const handleOnEdit = useCallback(
    (id: string, lastRunId: number) => {
      navigate(`/testsuite/${id}/run/${lastRunId}`);
    },
    [navigate]
  );

  return (
    <>
      <S.Wrapper>
        <S.HeaderContainer>
          <S.TitleText>All Test Suites</S.TitleText>
          <VariableSetSelector />
        </S.HeaderContainer>

        <S.ActionsContainer>
          <HomeFilters
            onSearch={pagination.search}
            onSortBy={(sortBy, sortDirection) => setParameters({sortBy, sortDirection})}
            isEmpty={pagination.list?.length === 0}
          />
          <AllowButton
            operation={Operation.Edit}
            ButtonComponent={CreateButton}
            data-cy="create-button"
            onClick={() => setIsCreateTestSuiteOpen(true)}
            type="primary"
          >
            Create
          </AllowButton>
        </S.ActionsContainer>

        <Pagination<TestSuite>
          emptyComponent={
            !testListData?.total ? (
              <Empty
                title="No Test Suites to Display... Yet!"
                message="To set up your test suits and experience the interconnected testing magic, let's kickstart by creating your first test. Ready to boost your test coverage and efficiency? Let's dive in!"
                action={
                  <S.Button onClick={() => navigate('/')} type="primary">
                    Go To Tests Page
                  </S.Button>
                }
              />
            ) : (
              <Empty
                title="No Test Suits in Sight!"
                message={
                  <>
                    It looks a bit empty here, doesn&apos;t it? Let&apos;s start by adding your first test suites. If
                    you want to learn more about test suits just click{' '}
                    <S.Link href={ADD_TEST_SUITE_URL} target="_blank">
                      here
                    </S.Link>
                    .
                  </>
                }
                action={
                  <AllowButton
                    operation={Operation.Edit}
                    ButtonComponent={CreateButton}
                    onClick={() => setIsCreateTestSuiteOpen(true)}
                    type="primary"
                  >
                    Create Your First Test Suite
                  </AllowButton>
                }
              />
            )
          }
          loadingComponent={<Loading />}
          {...pagination}
        >
          <S.TestListContainer data-cy="test-list">
            {pagination.list?.map(suite => (
              <TestSuiteCard
                key={suite.id}
                onEdit={handleOnEdit}
                onDelete={onDelete}
                onRun={handleOnRun}
                onViewAll={handleOnViewAll}
                testSuite={suite}
              />
            ))}
          </S.TestListContainer>
        </Pagination>
      </S.Wrapper>

      <CreateTestSuiteModal isOpen={isCreateTestSuiteOpen} onClose={() => setIsCreateTestSuiteOpen(false)} />
    </>
  );
};

export default Resources;
