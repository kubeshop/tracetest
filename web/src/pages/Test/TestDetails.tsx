import {Button, Typography} from 'antd';
import {useCallback} from 'react';

import {Steps} from 'components/GuidedTour/testDetailsStepList';
import InfiniteScroll from 'components/InfiniteScroll';
import ResultCardList from 'components/RunCardList';
import SearchInput from 'components/SearchInput';
import useInfiniteScroll from 'hooks/useInfiniteScroll';
import {useGetRunListQuery, useRunTestMutation} from 'redux/apis/TraceTest.api';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {TTestRun} from 'types/TestRun.types';
import * as S from './Test.styled';

const {onRunTest} = TestAnalyticsService;

interface IProps {
  onSelectResult: (result: TTestRun) => void;
  testId: string;
}

const TestDetails = ({onSelectResult, testId}: IProps) => {
  const [runTest, result] = useRunTestMutation();
  const {
    list: resultList,
    hasMore,
    loadMore,
    isLoading,
  } = useInfiniteScroll<TTestRun, {testId: string}>(useGetRunListQuery, {
    testId,
  });

  const handleRunTest = useCallback(async () => {
    if (testId) {
      onRunTest();
      const testResult = await runTest({testId}).unwrap();
      onSelectResult(testResult);
    }
  }, [onSelectResult, runTest, testId]);

  return (
    <>
      <S.TestDetailsHeader>
        <SearchInput
          onSearch={() => {
            // eslint-disable-next-line no-console
            console.log('onSearch');
          }}
          placeholder="Search test result (Not implemented yet)"
        />
        <Button
          onClick={handleRunTest}
          loading={result.isLoading}
          type="primary"
          data-cy="test-details-run-test-button"
          ghost
          data-tour={GuidedTourService.getStep(GuidedTours.TestDetails, Steps.RunTest)}
        >
          Run Test
        </Button>
      </S.TestDetailsHeader>

      <InfiniteScroll
        loadMore={loadMore}
        isLoading={isLoading}
        hasMore={hasMore}
        shouldTrigger={Boolean(resultList.length)}
        emptyComponent={
          <S.EmptyStateContainer>
            <S.EmptyStateIcon />
            <Typography.Text disabled>No Runs</Typography.Text>
          </S.EmptyStateContainer>
        }
      >
        <ResultCardList testId={testId} resultList={resultList} />
      </InfiniteScroll>
    </>
  );
};

export default TestDetails;
