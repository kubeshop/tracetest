import CreateTestModal from 'components/CreateTestModal/CreateTestModal';
import CreateTransactionModal from 'components/CreateTransactionModal/CreateTransactionModal';
import Empty from 'components/Empty';
import Pagination from 'components/Pagination';
import TestCard from 'components/ResourceCard/TestCard';
import TransactionCard from 'components/ResourceCard/TransactionCard';
import {SortBy, SortDirection, sortOptions} from 'constants/Test.constants';
import useDeleteResource from 'hooks/useDeleteResource';
import usePagination from 'hooks/usePagination';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {useCallback, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {useGetResourcesQuery} from 'redux/apis/TraceTest.api';
import {ADD_TEST_SPECS_DOCUMENTATION_URL} from 'constants/Common.constants';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import {ResourceType} from 'types/Resource.type';
import useTransactionCrud from 'providers/Transaction/hooks/useTransactionCrud';
import Resource from 'models/Resource.model';
import Transaction from 'models/Transaction.model';
import Test from 'models/Test.model';
import * as S from './Home.styled';
import HomeActions from './HomeActions';
import HomeFilters from './HomeFilters';
import Loading from './Loading';

const {onTestClick} = HomeAnalyticsService;
type TParameters = {sortBy: SortBy; sortDirection: SortDirection};
const [{params: defaultSort}] = sortOptions;

const Resources = () => {
  const [isCreateTransactionOpen, setIsCreateTransactionOpen] = useState(false);
  const [isCreateTestOpen, setIsCreateTestOpen] = useState(false);
  const [parameters, setParameters] = useState<TParameters>(defaultSort);

  const pagination = usePagination<Resource, TParameters>(useGetResourcesQuery, parameters);
  const onDeleteResource = useDeleteResource();
  const {runTest} = useTestCrud();
  const {runTransaction} = useTransactionCrud();
  const navigate = useNavigate();

  const handleOnRun = useCallback(
    (resource: Transaction | Test, type: ResourceType) => {
      if (type === ResourceType.Test) runTest(resource as Test);
      else if (type === ResourceType.Transaction) runTransaction(resource as Transaction);
    },
    [runTest, runTransaction]
  );

  const handleOnViewAll = useCallback((id: string) => {
    onTestClick(id);
  }, []);

  const handleOnEdit = useCallback(
    (id: string, lastRunId: number, type: ResourceType) => {
      if (type === ResourceType.Test) navigate(`/test/${id}/run/${lastRunId}`);
      else if (type === ResourceType.Transaction) navigate(`/transaction/${id}/run/${lastRunId}`);
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
            onCreateTransaction={() => setIsCreateTransactionOpen(true)}
            onCreateTest={() => setIsCreateTestOpen(true)}
          />
        </S.ActionsContainer>

        <Pagination<Resource>
          emptyComponent={
            <Empty
              title="You have not created any tests yet"
              message={
                <>
                  Use the Create button to create your first test. Learn more about test or transactions{' '}
                  <a href={ADD_TEST_SPECS_DOCUMENTATION_URL} target="_blank">
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
                <TransactionCard
                  key={resource.item.id}
                  onEdit={handleOnEdit}
                  onDelete={onDeleteResource}
                  onRun={handleOnRun}
                  onViewAll={handleOnViewAll}
                  transaction={resource.item as Transaction}
                />
              )
            )}
          </S.TestListContainer>
        </Pagination>
      </S.Wrapper>

      <CreateTestModal isOpen={isCreateTestOpen} onClose={() => setIsCreateTestOpen(false)} />
      <CreateTransactionModal isOpen={isCreateTransactionOpen} onClose={() => setIsCreateTransactionOpen(false)} />
    </>
  );
};

export default Resources;
