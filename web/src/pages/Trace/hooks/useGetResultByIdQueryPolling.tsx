import usePolling from '../../../hooks/usePolling';
import {TestState} from '../../../constants/TestRunResult.constants';

export function useGetResultByIdQueryPolling(refetchTrace: any, isError: boolean, testResultDetails: any): void {
  usePolling({
    callback: refetchTrace,
    delay: 1000,
    isPolling:
      isError ||
      testResultDetails?.state === TestState.AWAITING_TRACE ||
      testResultDetails?.state === TestState.EXECUTING,
  });
}
