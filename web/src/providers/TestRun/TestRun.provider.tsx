import {noop} from 'lodash';
import {createContext, useCallback, useContext, useEffect, useMemo, useState} from 'react';

import TestRun, {isRunStateFinished} from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import TracetestAPI from 'redux/apis/Tracetest';
import TestProvider from '../Test';

const {useGetRunByIdQuery, useGetRunEventsQuery, useStopRunMutation} = TracetestAPI.instance;

interface IContext {
  run: TestRun;
  isError: boolean;
  isLoadingStop: boolean;
  runEvents: TestRunEvent[];
  stopRun(): void;
}

export const Context = createContext<IContext>({
  run: {} as TestRun,
  isError: false,
  isLoadingStop: false,
  runEvents: [],
  stopRun: noop,
});

interface IProps {
  testId: string;
  runId?: number;
  children: React.ReactNode;
}

export const useTestRun = () => useContext(Context);

const POLLING_INTERVAL = 5000;

const TestRunProvider = ({children, testId, runId = 0}: IProps) => {
  const [pollingInterval, setPollingInterval] = useState<number | undefined>(POLLING_INTERVAL);
  const {data: run, isError} = useGetRunByIdQuery({testId, runId}, {skip: !runId, pollingInterval});
  const {data: runEvents = []} = useGetRunEventsQuery({testId, runId}, {skip: !runId});
  const [stopRunAction, {isLoading: isLoadingStop}] = useStopRunMutation();

  const stopRun = useCallback(async () => {
    await stopRunAction({runId, testId});
  }, [runId, stopRunAction, testId]);

  const value = useMemo<IContext>(
    () => ({run: run!, isError, isLoadingStop, runEvents, stopRun}),
    [run, isError, isLoadingStop, runEvents, stopRun]
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
