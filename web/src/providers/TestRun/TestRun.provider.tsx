import {noop} from 'lodash';
import {createContext, useContext, useMemo} from 'react';
import {useGetRunByIdQuery} from 'redux/apis/TraceTest.api';
import {TTestRun} from 'types/TestRun.types';

interface IContext {
  run: TTestRun;
  refetch(): void;
  isError: boolean;
}

export const Context = createContext<IContext>({
  run: {} as TTestRun,
  refetch: noop,
  isError: false,
});

interface IProps {
  testId: string;
  runId?: string;
  children: React.ReactNode;
}

export const useTestRun = () => useContext(Context);

const TestRunProvider = ({children, testId, runId = ''}: IProps) => {
  const {data: run, refetch, isError} = useGetRunByIdQuery({testId, runId}, {skip: !runId});

  const value = useMemo<IContext>(() => ({run: run!, refetch, isError}), [run, refetch, isError]);

  return run ? <Context.Provider value={value}>{children}</Context.Provider> : <div data-cy="loading_test_run" />;
};

export default TestRunProvider;
