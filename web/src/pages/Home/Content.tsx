import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';

import InfiniteScroll from 'components/InfiniteScroll';
import SearchInput from 'components/SearchInput';
import TestCard from 'components/TestCard';
import TestCardV2 from 'components/TestCard/TestCardV2';
import useDeleteTest from 'hooks/useDeleteTest';
import useInfiniteScroll from 'hooks/useInfiniteScroll';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {useGetTestListQuery} from 'redux/apis/TraceTest.api';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import {TTest} from 'types/Test.types';
import ExperimentalFeature from 'utils/ExperimentalFeature';
import * as S from './Home.styled';
import HomeActions from './HomeActions';
import NoResults from './NoResults';
import Loading from './Loading';

const {onTestClick} = HomeAnalyticsService;

const Content = () => {
  const {hasMore, isLoading, list, loadMore, search} = useInfiniteScroll<TTest, {}>(useGetTestListQuery, {});
  const onDelete = useDeleteTest();
  const navigate = useNavigate();
  const {runTest} = useTestCrud();

  const handleOnViewAllClick = useCallback(
    (testId: string) => {
      onTestClick(testId);
      navigate(`/test/${testId}`);
    },
    [navigate]
  );

  const handleOnRun = useCallback(
    (testId: string) => {
      if (testId) runTest(testId);
    },
    [runTest]
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

      <InfiniteScroll
        loadMore={loadMore}
        isLoading={isLoading}
        hasMore={hasMore}
        shouldTrigger={Boolean(list.length)}
        emptyComponent={<NoResults />}
        loadingComponent={<Loading />}
      >
        <S.TestListContainer data-cy="test-list">
          {list?.map(test =>
            ExperimentalFeature.isEnabled('transactions') ? (
              <TestCardV2
                key={test.id}
                onDelete={onDelete}
                onRun={handleOnRun}
                onViewAll={handleOnViewAllClick}
                test={test}
              />
            ) : (
              <TestCard
                test={test}
                onClick={handleOnViewAllClick}
                onDelete={onDelete}
                onRunTest={handleOnRun}
                key={test.id}
              />
            )
          )}
        </S.TestListContainer>
      </InfiniteScroll>
    </S.Wrapper>
  );
};

export default Content;
