import {createContext, useContext, useMemo} from 'react';
import {useGetRunByIdQuery} from 'redux/apis/TraceTest.api';
import {TTestRun} from 'types/TestRun.types';
import TestProvider from '../Test';

interface IContext {
  run: TTestRun;
  isError: boolean;
}

export const Context = createContext<IContext>({
  run: {} as TTestRun,
  isError: false,
});

interface IProps {
  testId: string;
  runId?: string;
  children: React.ReactNode;
}

export const useTestRun = () => useContext(Context);

const TestRunProvider = ({children, testId, runId = ''}: IProps) => {
  const {data: run, isError} = useGetRunByIdQuery({testId, runId}, {skip: !runId});

  const value = useMemo<IContext>(() => ({run: run!, isError}), [run, isError]);

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
