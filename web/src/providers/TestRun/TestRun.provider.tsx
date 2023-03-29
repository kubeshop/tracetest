import {createContext, useContext, useMemo} from 'react';
import {useGetRunByIdQuery, useGetRunEventsQuery} from 'redux/apis/TraceTest.api';
import TestRun from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import TestProvider from '../Test';

interface IContext {
  run: TestRun;
  isError: boolean;
  runEvents: TestRunEvent[];
}

export const Context = createContext<IContext>({
  run: {} as TestRun,
  isError: false,
  runEvents: [],
});

interface IProps {
  testId: string;
  runId?: string;
  children: React.ReactNode;
}

export const useTestRun = () => useContext(Context);

const TestRunProvider = ({children, testId, runId = ''}: IProps) => {
  const {data: run, isError} = useGetRunByIdQuery({testId, runId}, {skip: !runId});
  const {data: runEvents = []} = useGetRunEventsQuery({testId, runId}, {skip: !runId});

  const value = useMemo<IContext>(() => ({run: run!, isError, runEvents}), [run, isError, runEvents]);

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
