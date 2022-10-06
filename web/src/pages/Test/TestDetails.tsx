import {Button, Typography} from 'antd';
import {useCallback} from 'react';

import {Steps} from 'components/GuidedTour/testDetailsStepList';
import Pagination from 'components/Pagination';
import ResultCardList from 'components/RunCardList';
import usePagination from 'hooks/usePagination';
import {useGetRunListQuery} from 'redux/apis/TraceTest.api';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {TTestRun} from 'types/TestRun.types';
import * as S from './Test.styled';
import useTestCrud from '../../providers/Test/hooks/useTestCrud';

interface IProps {
  testId: string;
}

const TestDetails = ({testId}: IProps) => {
  const {hasNext, hasPrev, isEmpty, isFetching, isLoading, list, loadNext, loadPrev} = usePagination<
    TTestRun,
    {testId: string}
  >(useGetRunListQuery, {
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

      <Pagination
        emptyComponent={
          <S.EmptyStateContainer>
            <S.EmptyStateIcon />
            <Typography.Text disabled>No Runs</Typography.Text>
          </S.EmptyStateContainer>
        }
        hasNext={hasNext}
        hasPrev={hasPrev}
        isEmpty={isEmpty}
        isFetching={isFetching}
        isLoading={isLoading}
        loadNext={loadNext}
        loadPrev={loadPrev}
      >
        <ResultCardList testId={testId} resultList={list} />
      </Pagination>
    </>
  );
};

export default TestDetails;
