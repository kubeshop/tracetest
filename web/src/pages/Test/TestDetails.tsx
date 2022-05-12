import {Button} from 'antd';
import {FC, useCallback} from 'react';
import {useGetResultListQuery, useRunTestMutation} from 'redux/apis/Test.api';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {Steps} from 'components/GuidedTour/testDetailsStepList';
import useGuidedTour from 'hooks/useGuidedTour';
import * as S from './Test.styled';
import TestAnalyticsService from '../../services/Analytics/TestAnalytics.service';
import {ITestRunResult} from '../../types/TestRunResult.types';
import ResultCardList from '../../components/ResultCardList';
import useInfiniteScroll from './hooks/useInfiniteScroll';
import InfiniteScroll from '../../components/InfiniteScroll';
import SearchInput from '../../components/SearchInput';

const {onRunTest} = TestAnalyticsService;

type TTestDetailsProps = {
  testId: string;
  onSelectResult: (result: ITestRunResult) => void;
};

const TestDetails: FC<TTestDetailsProps> = ({testId, onSelectResult}) => {
  const [runTest, result] = useRunTestMutation();
  const {
    list: resultList,
    hasMore,
    loadMore,
    isLoading,
  } = useInfiniteScroll<ITestRunResult, {testId: string}>(useGetResultListQuery, {
    testId,
  });

  useGuidedTour(GuidedTours.TestDetails);

  const handleRunTest = useCallback(async () => {
    if (testId) {
      onRunTest(testId);
      const testResult = await runTest(testId).unwrap();
      onSelectResult(testResult);
    }
  }, [onSelectResult, runTest, testId]);

  return (
    <>
      <S.TestDetailsHeader>
        <SearchInput onSearch={() => console.log('onSearch')} placeholder="Search test result (Not implemented yet)" />
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
      >
        <ResultCardList resultList={resultList} />
      </InfiniteScroll>
    </>
  );
};

export default TestDetails;
