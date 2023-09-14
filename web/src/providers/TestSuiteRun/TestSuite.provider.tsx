import {createContext, useContext, useEffect, useMemo, useState} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import TestSuiteRun, {isRunStateFinished} from 'models/TestSuiteRun.model';
import TestSuiteProvider from '../TestSuite/TestSuite.provider';

const {useGetTestSuiteRunByIdQuery} = TracetestAPI.instance;

interface IContext {
  run: TestSuiteRun;
}

export const Context = createContext<IContext>({
  run: {} as TestSuiteRun,
});

interface IProps {
  testSuiteId: string;
  runId: number;
  children: React.ReactNode;
}

export const useTestSuiteRun = () => useContext(Context);

const POLLING_INTERVAL = 5000;

const TestSuiteRunProvider = ({children, testSuiteId, runId}: IProps) => {
  const [pollingInterval, setPollingInterval] = useState<number | undefined>(POLLING_INTERVAL);
  const {data: run} = useGetTestSuiteRunByIdQuery({testSuiteId, runId}, {pollingInterval});
  const value = useMemo<IContext>(() => ({run: run!}), [run]);

  useEffect(() => {
    const shouldStopPolling = run?.state && isRunStateFinished(run.state);
    setPollingInterval(shouldStopPolling ? undefined : POLLING_INTERVAL);
  }, [run?.state]);

  return run ? (
    <TestSuiteProvider testSuiteId={testSuiteId} version={run.version}>
      <Context.Provider value={value}>{children}</Context.Provider>
    </TestSuiteProvider>
  ) : null;
};

export default TestSuiteRunProvider;
