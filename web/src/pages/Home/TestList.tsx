import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';

import InfiniteScroll from 'components/InfiniteScroll';
import TestCard from 'components/TestCard';
import useInfiniteScroll from 'hooks/useInfiniteScroll';
import {useGetTestListQuery, useRunTestMutation} from 'redux/apis/TraceTest.api';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {TTest} from 'types/Test.types';
import * as S from './Home.styled';
import NoResults from './NoResults';
import {useMenuDeleteCallback} from './useMenuDeleteCallback';

const {onTestClick} = HomeAnalyticsService;

const TestList = () => {
  const navigate = useNavigate();
  const [runTest] = useRunTestMutation();
  const {list: resultList, hasMore, loadMore, isLoading} = useInfiniteScroll<TTest, {}>(useGetTestListQuery, {});

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

  const onDelete = useMenuDeleteCallback();

  return (
    <InfiniteScroll
      loadMore={loadMore}
      isLoading={isLoading}
      hasMore={hasMore}
      shouldTrigger={Boolean(resultList.length)}
      emptyComponent={<NoResults />}
    >
      <S.TestListContainer data-cy="test-list">
        {resultList?.map(test => (
          <TestCard test={test} onClick={onClick} onDelete={onDelete} onRunTest={onRunTest} key={test.id} />
        ))}
      </S.TestListContainer>
    </InfiniteScroll>
  );
};

export default TestList;
