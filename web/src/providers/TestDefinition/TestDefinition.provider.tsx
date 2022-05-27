import {noop} from 'lodash';
import {createContext, useContext, useMemo, useEffect} from 'react';
import {useGetTestByIdQuery} from '../../redux/apis/TraceTest.api';
import {useAppSelector} from '../../redux/hooks';
import TestDefinitionSelectors from '../../selectors/TestDefinition.selectors';
import {TAssertionResults} from '../../types/Assertion.types';
import {TTest} from '../../types/Test.types';
import {TTestDefinitionEntry} from '../../types/TestDefinition.types';
import {useTestRun} from '../TestRun/TestRun.provider';
import useTestDefinitionCrud from './hooks/useTestDefinitionCrud';

interface IContext {
  add(testDefinition: TTestDefinitionEntry): void;
  update(selector: string, testDefinition: TTestDefinitionEntry): void;
  remove(selector: string): void;
  publish(): void;
  cancel(): void;
  dryRun(definitionList: TTestDefinitionEntry[]): void;
  assertionResults?: TAssertionResults;
  definitionList: TTestDefinitionEntry[];
  isLoading: boolean;
  isError: boolean;
  isDraftMode: boolean;
  test?: TTest;
}

export const Context = createContext<IContext>({
  add: noop,
  update: noop,
  remove: noop,
  publish: noop,
  dryRun: noop,
  cancel: noop,
  isLoading: false,
  isError: false,
  isDraftMode: false,
  definitionList: [],
});

interface IProps {
  testId: string;
  runId: string;
}

export const useTestDefinition = () => useContext(Context);

const TestDefinitionProvider: React.FC<IProps> = ({children, testId, runId}) => {
  const {run} = useTestRun();
  const assertionResults = useAppSelector(state => TestDefinitionSelectors.selectAssertionResults(state));
  const definitionList = useAppSelector(state => TestDefinitionSelectors.selectDefinitionList(state));
  const isLoading = useAppSelector(state => TestDefinitionSelectors.selectIsLoading(state));
  const isInitialized = useAppSelector(state => TestDefinitionSelectors.selectIsInitialized(state));
  const {data: test} = useGetTestByIdQuery({testId});

  const {add, cancel, publish, remove, dryRun, update, isDraftMode, init, reset} = useTestDefinitionCrud({
    testId,
    runId,
  });

  useEffect(() => {
    init(run.result);
  }, [init, isInitialized, reset, run.result]);

  useEffect(() => {
    return () => {
      reset();
    };
  }, [reset]);

  useEffect(() => {
    if (isInitialized && run.state === 'FINISHED') dryRun(definitionList);
  }, [dryRun, definitionList, isInitialized, run.state]);

  const value = useMemo<IContext>(
    () => ({
      add,
      remove,
      update,
      isLoading,
      isError: false,
      isDraftMode,
      publish,
      dryRun,
      assertionResults,
      definitionList,
      cancel,
      test,
    }),
    [add, assertionResults, cancel, definitionList, dryRun, isDraftMode, isLoading, publish, remove, update, test]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TestDefinitionProvider;
