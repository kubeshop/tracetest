import {Button, Typography} from 'antd';
import {useCallback} from 'react';

import {Steps} from 'components/GuidedTour/testDetailsStepList';
import InfiniteScroll from 'components/InfiniteScroll';
import ResultCardList from 'components/RunCardList';
import useInfiniteScroll from 'hooks/useInfiniteScroll';
import {useGetRunListQuery} from 'redux/apis/TraceTest.api';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {TTestRun} from 'types/TestRun.types';
import * as S from './Test.styled';
import useTestCrud from '../../providers/Test/hooks/useTestCrud';

interface IProps {
  testId: string;
}

const TestDetails = ({testId}: IProps) => {
  const {
    list: resultList,
    hasMore,
    loadMore,
    isLoading,
  } = useInfiniteScroll<TTestRun, {testId: string}>(useGetRunListQuery, {
    testId,
  });

  const {runTest, isLoadingRunTest} = useTestCrud();

  const handleRunTest = useCallback(async () => {
    if (testId) runTest(testId);
  }, [runTest, testId]);

  return (
    <>
      <S.TestDetailsHeader>
        <div />
        <Button
          onClick={handleRunTest}
          loading={isLoadingRunTest}
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
