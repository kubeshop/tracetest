import React, {useState} from 'react';
import {useGetResultByIdQuery, useGetTestByIdQuery} from 'redux/apis/Test.api';

import {GuidedTours} from 'services/GuidedTour.service';
import useGuidedTour from 'hooks/useGuidedTour';
import {ITestRunResult} from 'types/TestRunResult.types';
import {useAppDispatch} from 'redux/hooks';
import {TestState} from 'constants/TestRunResult.constants';
import {TraceError} from 'components/Trace/TraceError/TraceError';
import {useGetResultByIdQueryPolling} from 'components/Trace/useGetResultByIdQueryPolling';
import {TraceComponent} from 'components/Trace/TraceComponent/TraceComponent';
import {useUpdateTestResultEffect} from 'components/Trace/useUpdateTestResultEffect';
import {useUpdateFirstLoadEffect} from 'components/Trace/useUpdateFirstLoadEffect';
import {visiblePortionFuction} from 'components/Trace/visiblePortionFuction';
import 'react-reflex/styles.css';

interface TraceProps {
  testId: string;
  testResultId: string;

  onDismissTrace(): void;

  onRunTest(result: ITestRunResult): void;
}

export const Trace: React.FC<TraceProps> = ({testId, testResultId, onDismissTrace, onRunTest}): JSX.Element => {
  const dispatch = useAppDispatch();
  const {visiblePortion, height} = visiblePortionFuction();
  const [isFirstLoad, setIsFirstLoad] = useState(true);
  const {data: test} = useGetTestByIdQuery(testId);
  useGuidedTour(GuidedTours.Trace);

  const query = useGetResultByIdQuery({testId, resultId: testResultId});
  const displayError = query.isError || query.data?.state === TestState.FAILED;
  useGetResultByIdQueryPolling(query.refetch, query.isError, query.data);
  useUpdateTestResultEffect(query.data, test, isFirstLoad, dispatch, testResultId);
  useUpdateFirstLoadEffect(query.data, test, setIsFirstLoad, dispatch, testResultId, testId);
  return (
    <>
      <TraceError onRunTest={onRunTest} testId={testId} displayError={displayError} onDismissTrace={onDismissTrace} />
      <TraceComponent
        displayError={displayError}
        minHeight={height}
        testResultDetails={query.data}
        test={test}
        visiblePortion={visiblePortion}
      />
    </>
  );
};
