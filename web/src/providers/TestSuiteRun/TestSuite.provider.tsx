import {createContext, useContext, useMemo} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import TestSuiteRun from 'models/TestSuiteRun.model';
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

const TestSuiteRunProvider = ({children, testSuiteId, runId}: IProps) => {
  const {data: run} = useGetTestSuiteRunByIdQuery({testSuiteId, runId});
  const value = useMemo<IContext>(() => ({run: run!}), [run]);

  return run ? (
    <TestSuiteProvider testSuiteId={testSuiteId} version={run.version}>
      <Context.Provider value={value}>{children}</Context.Provider>
    </TestSuiteProvider>
  ) : null;
};

export default TestSuiteRunProvider;
