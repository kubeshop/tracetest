import {useEffect} from 'react';
import {ITestRunResult} from 'types/TestRunResult.types';
import {ITest} from 'types/Test.types';
import {AppDispatch} from 'redux/store';
import {updateTestResult} from 'redux/slices/ResultList.slice';

export function useUpdateTestResultEffect(
  testResultDetails: ITestRunResult | undefined,
  test: ITest | undefined,
  isFirstLoad: boolean,
  dispatch: AppDispatch,
  testResultId: string
) {
  useEffect(() => {
    if (testResultDetails && test && !isFirstLoad) {
      dispatch(
        updateTestResult({
          trace: testResultDetails.trace!,
          resultId: testResultId,
          test,
        })
      );
    }
  }, [test, dispatch, isFirstLoad, testResultDetails, testResultId]);
}
