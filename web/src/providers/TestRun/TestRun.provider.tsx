import {noop} from 'lodash';
import {createContext, useContext, useMemo} from 'react';
import {useGetRunByIdQuery} from '../../redux/apis/TraceTest.api';
import {TTestRun} from '../../types/TestRun.types';

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
  runId: string;
}

export const useTestRun = () => useContext(Context);

const TestRunProvider: React.FC<IProps> = ({children, testId, runId}) => {
  const {data: run, refetch, isError} = useGetRunByIdQuery({testId, runId});

  const value = useMemo(() => ({run, refetch, isError}), [refetch, run, isError]) as IContext;

  return run ? <Context.Provider value={value}>{children}</Context.Provider> : <div data-cy="not_initialized" />;
};

export default TestRunProvider;
