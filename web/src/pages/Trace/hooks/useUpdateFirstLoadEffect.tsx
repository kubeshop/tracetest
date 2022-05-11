import {useEffect} from 'react';
import {isEmpty} from 'lodash';
import {ITest} from 'types/Test.types';
import {AppDispatch} from 'redux/store';
import {replace, updateTestResult} from 'redux/slices/ResultList.slice';
import {ITestRunResult} from 'types/TestRunResult.types';

export function useUpdateFirstLoadEffect(
  testResultDetails: ITestRunResult | undefined,
  test: ITest | undefined,
  setIsFirstLoad: (value: ((prevState: boolean) => boolean) | boolean) => void,
  dispatch: AppDispatch,
  testResultId: string,
  testId: string
) {
  useEffect(() => {
    if (testResultDetails && !isEmpty(testResultDetails.trace) && !testResultDetails?.assertionResult && test) {
      setIsFirstLoad(false);
      dispatch(
        updateTestResult({
          trace: testResultDetails.trace!,
          resultId: testResultId,
          test: test!,
        })
      );
    } else if (testResultDetails?.assertionResult && test) {
      setIsFirstLoad(false);

      dispatch(
        // @ts-ignore
        replace({
          resultId: testResultId,
          assertionResult: testResultDetails?.assertionResult!,
          test,
          trace: testResultDetails?.trace!,
        })
      );
    }
  }, [setIsFirstLoad, testResultDetails, test, testResultId, testId, dispatch]);
}
