import CreateTestModal from 'components/CreateTestModal/CreateTestModal';
import CreateTestSuiteModal from 'components/CreateTestSuiteModal/CreateTestSuiteModal';
import Empty from 'components/Empty';
import Pagination from 'components/Pagination';
import TestCard from 'components/ResourceCard/TestCard';
import TestSuiteCard from 'components/ResourceCard/TestSuiteCard';
import {SortBy, SortDirection, sortOptions} from 'constants/Test.constants';
import useDeleteResource from 'hooks/useDeleteResource';
import usePagination from 'hooks/usePagination';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {useCallback, useState} from 'react';
import {useGetResourcesQuery} from 'redux/apis/Tracetest';
import {ADD_TEST_URL} from 'constants/Common.constants';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import {ResourceType} from 'types/Resource.type';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import useTestSuiteCrud from 'providers/TestSuite/hooks/useTestSuiteCrud';
import Resource from 'models/Resource.model';
import TestSuite from 'models/TestSuite.model';
import Test from 'models/Test.model';
import * as S from './Home.styled';
import HomeActions from './HomeActions';
import HomeFilters from './HomeFilters';
import Loading from './Loading';

const {onTestClick} = HomeAnalyticsService;
type TParameters = {sortBy: SortBy; sortDirection: SortDirection};
const [{params: defaultSort}] = sortOptions;

const Resources = () => {
  const [isCreateTestSuiteOpen, setIsCreateTestSuiteOpen] = useState(false);
  const [isCreateTestOpen, setIsCreateTestOpen] = useState(false);
  const [parameters, setParameters] = useState<TParameters>(defaultSort);

  const pagination = usePagination<Resource, TParameters>(useGetResourcesQuery, parameters);
  const onDeleteResource = useDeleteResource();
  const {runTest} = useTestCrud();
  const {runTestSuite} = useTestSuiteCrud();
  const {navigate} = useDashboard();

  const handleOnRun = useCallback(
    (resource: TestSuite | Test, type: ResourceType) => {
      if (type === ResourceType.Test) runTest({test: resource as Test});
      else if (type === ResourceType.TestSuite) runTestSuite(resource as TestSuite);
    },
    [runTest, runTestSuite]
  );

  const handleOnViewAll = useCallback((id: string) => {
    onTestClick(id);
  }, []);

  const handleOnEdit = useCallback(
    (id: string, lastRunId: number, type: ResourceType) => {
      if (type === ResourceType.Test) navigate(`/test/${id}/run/${lastRunId}`);
      else if (type === ResourceType.TestSuite) navigate(`/testsuite/${id}/run/${lastRunId}`);
    },
    [navigate]
  );

  return (
    <>
      <S.Wrapper>
        <S.HeaderContainer>
          <S.TitleText>All Tests</S.TitleText>
        </S.HeaderContainer>

        <S.ActionsContainer>
          <HomeFilters
            onSearch={pagination.search}
            onSortBy={(sortBy, sortDirection) => setParameters({sortBy, sortDirection})}
            isEmpty={pagination.list?.length === 0}
          />
          <HomeActions
            onCreateTestSuite={() => setIsCreateTestSuiteOpen(true)}
            onCreateTest={() => setIsCreateTestOpen(true)}
          />
        </S.ActionsContainer>

        <Pagination<Resource>
          emptyComponent={
            <Empty
              title="You have not created any tests yet"
              message={
                <>
                  Use the Create button to create your first test. Learn more about test or test suites{' '}
                  <a href={ADD_TEST_URL} target="_blank">
                    here.
                  </a>
                </>
              }
            />
          }
          loadingComponent={<Loading />}
          {...pagination}
        >
          <S.TestListContainer data-cy="test-list">
            {pagination.list?.map(resource =>
              resource.type === ResourceType.Test ? (
                <TestCard
                  key={resource.item.id}
                  onEdit={handleOnEdit}
                  onDelete={onDeleteResource}
                  onRun={handleOnRun}
                  onViewAll={handleOnViewAll}
                  test={resource.item as Test}
                />
              ) : (
                <TestSuiteCard
                  key={resource.item.id}
                  onEdit={handleOnEdit}
                  onDelete={onDeleteResource}
                  onRun={handleOnRun}
                  onViewAll={handleOnViewAll}
                  testSuite={resource.item as TestSuite}
                />
              )
            )}
          </S.TestListContainer>
        </Pagination>
      </S.Wrapper>

      <CreateTestModal isOpen={isCreateTestOpen} onClose={() => setIsCreateTestOpen(false)} />
      <CreateTestSuiteModal isOpen={isCreateTestSuiteOpen} onClose={() => setIsCreateTestSuiteOpen(false)} />
    </>
  );
};

export default Resources;
