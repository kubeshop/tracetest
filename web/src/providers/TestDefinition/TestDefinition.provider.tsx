import {noop} from 'lodash';

import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {createContext, useCallback, useContext, useEffect, useMemo} from 'react';
import {useGetTestByIdQuery} from 'redux/apis/TraceTest.api';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {
  setSelectedAssertion as setSelectedAssertionAction,
  setViewResultsMode,
} from 'redux/slices/TestDefinition.slice';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TAssertionResultEntry, TAssertionResults} from 'types/Assertion.types';
import {TTest} from 'types/Test.types';
import {TTestDefinitionEntry} from 'types/TestDefinition.types';
import {ResultViewModes} from 'constants/Test.constants';
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
  viewResultsMode: ResultViewModes;
  setSelectedAssertion(assertionResult?: TAssertionResultEntry): void;
  changeViewResultsMode(viewResultsMode: ResultViewModes): void;
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
  changeViewResultsMode: noop,
  isLoading: false,
  isError: false,
  isDraftMode: false,
  definitionList: [],
  viewResultsMode: ResultViewModes.Wizard,
  setSelectedAssertion: noop,
});

interface IProps {
  testId: string;
  runId: string;
  children: React.ReactNode;
}

export const useTestDefinition = () => useContext(Context);

const TestDefinitionProvider = ({children, testId, runId}: IProps) => {
  const dispatch = useAppDispatch();
  const {run} = useTestRun();
  const assertionResults = useAppSelector(state => TestDefinitionSelectors.selectAssertionResults(state));
  const definitionList = useAppSelector(state => TestDefinitionSelectors.selectDefinitionList(state));
  const isDraftMode = useAppSelector(state => TestDefinitionSelectors.selectIsDraftMode(state));
  const isLoading = useAppSelector(state => TestDefinitionSelectors.selectIsLoading(state));
  const isInitialized = useAppSelector(state => TestDefinitionSelectors.selectIsInitialized(state));
  const viewResultsMode = useAppSelector(state => TestDefinitionSelectors.selectViewResultsMode(state));
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

  const setSelectedAssertion = useCallback(
    (assertionResult?: TAssertionResultEntry) => {
      dispatch(setSelectedAssertionAction(assertionResult));
    },
    [dispatch]
  );

  const changeViewResultsMode = useCallback(
    (mode: ResultViewModes) => {
      dispatch(setViewResultsMode(mode));
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
      viewResultsMode,
      cancel,
      test,
      setSelectedAssertion,
      revert,
      changeViewResultsMode,
    }),
    [
      add,
      remove,
      update,
      isLoading,
      isDraftMode,
      publish,
      runTest,
      dryRun,
      assertionResults,
      definitionList,
      viewResultsMode,
      cancel,
      test,
      setSelectedAssertion,
      revert,
      changeViewResultsMode,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TestDefinitionProvider;
