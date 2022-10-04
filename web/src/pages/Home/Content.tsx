import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';

import Pagination from 'components/Pagination';
import SearchInput from 'components/SearchInput';
import TestCard from 'components/TestCard';
import useDeleteTest from 'hooks/useDeleteTest';
import usePagination from 'hooks/usePagination';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {useGetTestListQuery} from 'redux/apis/TraceTest.api';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import {TTest} from 'types/Test.types';
import * as S from './Home.styled';
import HomeActions from './HomeActions';
import NoResults from './NoResults';
import Loading from './Loading';

const {onTestClick} = HomeAnalyticsService;

const Content = () => {
  const {hasNext, hasPrev, isEmpty, isFetching, isLoading, list, loadNext, loadPrev, search} = usePagination<TTest, {}>(
    useGetTestListQuery,
    {}
  );
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
    <S.Wrapper>
      <S.HeaderContainer>
        <S.TitleText>All Tests</S.TitleText>
      </S.HeaderContainer>

      <S.ActionsContainer>
        <SearchInput onSearch={value => search(value)} placeholder="Search test" />
        <HomeActions />
      </S.ActionsContainer>

      <Pagination
        emptyComponent={<NoResults />}
        hasNext={hasNext}
        hasPrev={hasPrev}
        isEmpty={isEmpty}
        isFetching={isFetching}
        isLoading={isLoading}
        loadingComponent={<Loading />}
        loadNext={loadNext}
        loadPrev={loadPrev}
      >
        <S.TestListContainer data-cy="test-list">
          {list?.map(test => (
            <TestCard key={test.id} onDelete={onDelete} onRun={handleOnRun} onViewAll={handleOnViewAll} test={test} />
          ))}
        </S.TestListContainer>
      </Pagination>
    </S.Wrapper>
  );
};

export default Content;
