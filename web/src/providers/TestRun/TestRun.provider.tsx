import {noop} from 'lodash';
import {createContext, useCallback, useContext, useEffect, useMemo, useState} from 'react';

import TestRun, {isRunStateFinished} from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import TracetestAPI from 'redux/apis/Tracetest';
import TestProvider from '../Test';

const {useGetRunByIdQuery, useGetRunEventsQuery, useStopRunMutation, useSkipPollingMutation} = TracetestAPI.instance;

interface IContext {
  run: TestRun;
  isError: boolean;
  isLoadingStop: boolean;
  isLoadingSkipPolling: boolean;
  runEvents: TestRunEvent[];
  stopRun(): void;
  skipPolling(): void;
}

export const Context = createContext<IContext>({
  run: {} as TestRun,
  isError: false,
  isLoadingStop: false,
  isLoadingSkipPolling: false,
  runEvents: [],
  stopRun: noop,
  skipPolling: noop,
});

interface IProps {
  testId: string;
  runId?: number;
  children: React.ReactNode;
}

export const useTestRun = () => useContext(Context);

const POLLING_INTERVAL = 1000;

const TestRunProvider = ({children, testId, runId = 0}: IProps) => {
  const [pollingInterval, setPollingInterval] = useState<number | undefined>(POLLING_INTERVAL);
  const {data: run, isError} = useGetRunByIdQuery({testId, runId}, {skip: !runId, pollingInterval});
  const {data: runEvents = []} = useGetRunEventsQuery({testId, runId}, {skip: !runId, pollingInterval});
  const [stopRunAction, {isLoading: isLoadingStop}] = useStopRunMutation();
  const [skipPollingAction, {isLoading: isLoadingSkipPolling}] = useSkipPollingMutation();

  const stopRun = useCallback(() => stopRunAction({runId, testId}), [runId, stopRunAction, testId]);
  const skipPolling = useCallback(() => skipPollingAction({runId, testId}), [runId, skipPollingAction, testId]);

  const value = useMemo<IContext>(
    () => ({run: run!, isError, isLoadingStop, isLoadingSkipPolling, runEvents, stopRun, skipPolling}),
    [run, isError, isLoadingStop, isLoadingSkipPolling, runEvents, stopRun, skipPolling]
  );

  useEffect(() => {
    const shouldStopPolling = run?.state && isRunStateFinished(run.state);
    setPollingInterval(shouldStopPolling ? undefined : POLLING_INTERVAL);
  }, [run?.state]);

  return run ? (
    <Context.Provider value={value}>
      <TestProvider testId={testId} version={run.testVersion}>
        {children}
      </TestProvider>
    </Context.Provider>
  ) : (
    <div data-cy="loading_test_run" />
  );
};

export default TestRunProvider;
