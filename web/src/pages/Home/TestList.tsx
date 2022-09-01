import InfiniteScroll from 'components/InfiniteScroll';
import TestCard from 'components/TestCard';
import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {useGetTestListQuery, useRunTestMutation} from 'redux/apis/TraceTest.api';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import useDeleteTest from 'hooks/useDeleteTest';
import useInfiniteScroll from 'hooks/useInfiniteScroll';
import {TTest} from 'types/Test.types';
import * as S from './Home.styled';
import NoResults from './NoResults';

const {onTestClick} = HomeAnalyticsService;

interface IProps {
  query: string;
}

const TestList = ({query}: IProps) => {
  const onDelete = useDeleteTest();
  const {list, isLoading, loadMore, hasMore} = useInfiniteScroll<TTest, {query: string}>(useGetTestListQuery, {query});
  const navigate = useNavigate();
  const [runTest] = useRunTestMutation();

  const onClick = useCallback(
    (testId: string) => {
      onTestClick(testId);
      navigate(`/test/${testId}`);
    },
    [navigate]
  );

  const onRunTest = useCallback(
    async (testId: string) => {
      if (testId) {
        TestAnalyticsService.onRunTest();
        const testRun = await runTest({testId}).unwrap();
        navigate(`/test/${testId}/run/${testRun.id}`);
      }
    },
    [navigate, runTest]
  );

  return (
    <InfiniteScroll
      loadMore={loadMore}
      isLoading={isLoading}
      hasMore={hasMore}
      shouldTrigger={Boolean(list.length)}
      emptyComponent={<NoResults />}
    >
      <S.TestListContainer data-cy="test-list">
        {list?.map(test => (
          <TestCard test={test} onClick={onClick} onDelete={onDelete} onRunTest={onRunTest} key={test.id} />
        ))}
      </S.TestListContainer>
    </InfiniteScroll>
  );
};

export default TestList;
