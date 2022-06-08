import {noop} from 'lodash';

import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {createContext, useCallback, useContext, useEffect, useMemo} from 'react';
import {useGetTestByIdQuery} from 'redux/apis/TraceTest.api';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {
  setAffectedSpans as setAffectedSpansAction,
  setSelectedAssertion as setSelectedAssertionAction,
} from 'redux/slices/TestDefinition.slice';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TAssertionResults} from 'types/Assertion.types';
import {TTest} from 'types/Test.types';
import {TTestDefinitionEntry} from 'types/TestDefinition.types';
import useTestDefinitionCrud from './hooks/useTestDefinitionCrud';

interface IContext {
  revert: (originalSelector: string) => void;
  add(testDefinition: TTestDefinitionEntry): void;
  update(selector: string, testDefinition: TTestDefinitionEntry): void;
  remove(selector: string): void;
  publish(): void;
  runTest(): void;
  cancel(): void;
  dryRun(definitionList: TTestDefinitionEntry[]): void;
  assertionResults?: TAssertionResults;
  definitionList: TTestDefinitionEntry[];
  isLoading: boolean;
  isError: boolean;
  isDraftMode: boolean;
  test?: TTest;
  setAffectedSpans(spanIds: string[]): void;
  setSelectedAssertion(selectorId: string): void;
}

export const Context = createContext<IContext>({
  add: noop,
  revert: () => noop,
  update: noop,
  remove: noop,
  publish: noop,
  runTest: noop,
  dryRun: noop,
  cancel: noop,
  isLoading: false,
  isError: false,
  isDraftMode: false,
  definitionList: [],
  setAffectedSpans: noop,
  setSelectedAssertion: noop,
});

interface IProps {
  testId: string;
  runId: string;
}

export const useTestDefinition = () => useContext(Context);

const TestDefinitionProvider: React.FC<IProps> = ({children, testId, runId}) => {
  const dispatch = useAppDispatch();
  const {run} = useTestRun();
  const assertionResults = useAppSelector(state => TestDefinitionSelectors.selectAssertionResults(state));
  const definitionList = useAppSelector(state => TestDefinitionSelectors.selectDefinitionList(state));
  const isDraftMode = useAppSelector(state => TestDefinitionSelectors.selectIsDraftMode(state));
  const isLoading = useAppSelector(state => TestDefinitionSelectors.selectIsLoading(state));
  const isInitialized = useAppSelector(state => TestDefinitionSelectors.selectIsInitialized(state));
  const {data: test} = useGetTestByIdQuery({testId});

  const {add, cancel, publish, runTest, remove, dryRun, update, init, reset, revert} = useTestDefinitionCrud({
    testId,
    runId,
    isDraftMode,
  });

  useEffect(() => {
    if (run.state === 'FINISHED') init(run.result);
  }, [init, run.result, run.state]);

  useEffect(() => {
    return () => {
      reset();
    };
  }, [reset]);

  useEffect(() => {
    if (isInitialized && run.state === 'FINISHED') dryRun(definitionList);
  }, [dryRun, definitionList, isInitialized, run.state]);

  const setAffectedSpans = useCallback(
    spanIds => {
      dispatch(setAffectedSpansAction(spanIds));
    },
    [dispatch]
  );

  const setSelectedAssertion = useCallback(
    selectorId => {
      dispatch(setSelectedAssertionAction(selectorId));
    },
    [dispatch]
  );

  const value = useMemo<IContext>(
    () => ({
      add,
      remove,
      update,
      isLoading,
      isError: false,
      isDraftMode,
      publish,
      runTest,
      dryRun,
      assertionResults,
      definitionList,
      cancel,
      test,
      setAffectedSpans,
      setSelectedAssertion,
      revert,
    }),
    [
      revert,
      add,
      assertionResults,
      cancel,
      definitionList,
      dryRun,
      isDraftMode,
      isLoading,
      publish,
      runTest,
      remove,
      update,
      test,
      setAffectedSpans,
      setSelectedAssertion,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TestDefinitionProvider;
