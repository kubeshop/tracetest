import {noop} from 'lodash';
import {createContext, useContext, useMemo, useEffect, useCallback} from 'react';

import {useTestRun} from 'providers/TestRun/TestRun.provider';
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
  setAffectedSpans(spanIds: string[]): void;
  setSelectedAssertion(selectorId: string): void;
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
  }, []);

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
      dryRun,
      assertionResults,
      definitionList,
      cancel,
      test,
      setAffectedSpans,
      setSelectedAssertion,
    }),
    [
      add,
      assertionResults,
      cancel,
      definitionList,
      dryRun,
      isDraftMode,
      isLoading,
      publish,
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
