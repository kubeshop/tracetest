import CreateTestModal from 'components/CreateTestModal/CreateTestModal';
import CreateTransactionModal from 'components/CreateTransactionModal/CreateTransactionModal';
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
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import {ResourceType, TResource} from 'types/Resource.type';
import {TTest} from 'types/Test.types';
import {TTransaction} from 'types/Transaction.types';
import useTransactionCrud from 'providers/Transaction/hooks/useTransactionCrud';
import * as S from './Home.styled';
import HomeActions from './HomeActions';
import HomeFilters from './HomeFilters';
import Loading from './Loading';
import NoResults from './NoResults';

const {onTestClick} = HomeAnalyticsService;
type TParameters = {sortBy: SortBy; sortDirection: SortDirection};
const [{params: defaultSort}] = sortOptions;

const Content = () => {
  const [isCreateTransactionOpen, setIsCreateTransactionOpen] = useState(false);
  const [isCreateTestOpen, setIsCreateTestOpen] = useState(false);
  const [parameters, setParameters] = useState<TParameters>(defaultSort);

  const navigate = useNavigate();
  const pagination = usePagination<TResource, TParameters>(useGetResourcesQuery, parameters);
  const onDeleteResource = useDeleteResource();
  const {runTest} = useTestCrud();
  const {runTransaction} = useTransactionCrud();

  const handleOnRun = useCallback(
    (id: string, type: ResourceType) => {
      if (type === ResourceType.Test) runTest(id);
      else if (type === ResourceType.Transaction) runTransaction(id);
    },
    [runTest, runTransaction]
  );

  const handleOnViewAll = useCallback(
    (id: string, type: ResourceType) => {
      onTestClick(id);
      navigate(`/${type}/${id}`);
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
          />
          <HomeActions
            onCreateTransaction={() => setIsCreateTransactionOpen(true)}
            onCreateTest={() => setIsCreateTestOpen(true)}
          />
        </S.ActionsContainer>

        <Pagination<TResource> emptyComponent={<NoResults />} loadingComponent={<Loading />} {...pagination}>
          <S.TestListContainer data-cy="test-list">
            {pagination.list?.map(resource =>
              resource.type === ResourceType.Test ? (
                <TestCard
                  key={resource.item.id}
                  onDelete={onDeleteResource}
                  onRun={handleOnRun}
                  onViewAll={handleOnViewAll}
                  test={resource.item as TTest}
                />
              ) : (
                <TransactionCard
                  key={resource.item.id}
                  onDelete={onDeleteResource}
                  onRun={handleOnRun}
                  onViewAll={handleOnViewAll}
                  transaction={resource.item as TTransaction}
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

export default Content;
