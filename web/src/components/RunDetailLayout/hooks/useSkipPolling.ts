import {useCallback, useEffect, useState} from 'react';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppDispatch} from 'redux/hooks';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import TestSpecsActions from 'redux/actions/TestSpecs.actions';
import {isRunStateFinished} from 'models/TestRun.model';
import {useTest} from 'providers/Test/Test.provider';

const useSkipPolling = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [shouldSave, setShouldSave] = useState(false);
  const {
    run: {id: runId, state},
    skipPolling,
  } = useTestRun();
  const {test} = useTest();

  const dispatch = useAppDispatch();
  const {navigate} = useDashboard();

  const onSkipPolling = useCallback(
    async (save: boolean) => {
      setIsLoading(true);
      skipPolling();
      setShouldSave(save);
    },
    [skipPolling]
  );

  const editAndReRun = useCallback(async () => {
    setShouldSave(false);
    const newRun = await dispatch(
      TestSpecsActions.publish({test: {...test, skipTraceCollection: true}, testId: test.id, runId})
    ).unwrap();
    setIsLoading(false);
    navigate(`/test/${test.id}/run/${newRun.id}`);
  }, [dispatch, navigate, runId, test]);

  useEffect(() => {
    if (isRunStateFinished(state) && shouldSave) editAndReRun();
    else if (isRunStateFinished(state)) setIsLoading(false);
  }, [dispatch, editAndReRun, shouldSave, state, test]);

  return {onSkipPolling, isLoading};
};

export default useSkipPolling;
