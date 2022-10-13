import CreateTestModal from 'components/CreateTestModal/CreateTestModal';
import CreateTransactionModal from 'components/CreateTransactionModal/CreateTransactionModal';
import Pagination from 'components/Pagination';
import TestCard from 'components/TestCard';
import {SortBy, SortDirection, sortOptions} from 'constants/Test.constants';
import useDeleteTest from 'hooks/useDeleteTest';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {useCallback, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import {useGetTestListQuery} from '../../redux/apis/TraceTest.api';
import {TTest} from '../../types/Test.types';
import * as S from './Home.styled';
import HomeActions from './HomeActions';
import HomeFilters from './HomeFilters';

const {onTestClick} = HomeAnalyticsService;
type TParameters = {sortBy: SortBy; sortDirection: SortDirection};
const [{params: defaultSort}] = sortOptions;

const Content = () => {
  const [isCreateTransactionOpen, setIsCreateTransactionOpen] = useState(false);
  const [isCreateTestOpen, setIsCreateTestOpen] = useState(false);

  const onDelete = useDeleteTest();
  const navigate = useNavigate();
  const {runTest} = useTestCrud();

  const handleOnRun = useCallback(
    (testId: string) => {
      if (testId) runTest(testId);
    },
    [runTest]
  );

  const handleOnViewAll = useCallback(
    (testId: string) => {
      onTestClick(testId);
      navigate(`/test/${testId}`);
    },
    [navigate]
  );

  return (
    <>
      <S.Wrapper>
        <S.HeaderContainer>
          <S.TitleText>All Tests</S.TitleText>
        </S.HeaderContainer>

        <Pagination<TTest, TParameters> query={useGetTestListQuery} defaultParameters={defaultSort}>
          {(pagination, [, setParameters]) => (
            <>
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
              <S.TestListContainer data-cy="test-list">
                {pagination.list?.map(test => (
                  <TestCard
                    key={test.id}
                    onDelete={onDelete}
                    onRun={handleOnRun}
                    onViewAll={handleOnViewAll}
                    test={test}
                  />
                ))}
              </S.TestListContainer>
            </>
          )}
        </Pagination>
      </S.Wrapper>
      <CreateTestModal isOpen={isCreateTestOpen} onClose={() => setIsCreateTestOpen(false)} />
      <CreateTransactionModal isOpen={isCreateTransactionOpen} onClose={() => setIsCreateTransactionOpen(false)} />
    </>
  );
};

export default Content;
