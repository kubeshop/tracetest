import {noop} from 'lodash';

import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {createContext, useCallback, useContext, useEffect, useMemo} from 'react';
import {useGetTestByIdQuery} from 'redux/apis/TraceTest.api';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TAssertionResultEntry, TAssertionResults} from 'types/Assertion.types';
import {TTest} from 'types/Test.types';
import {TTestSpecEntry} from 'types/TestSpecs.types';
import RouterActions from 'redux/actions/Router.actions';
import {RouterSearchFields} from 'constants/Common.constants';
import {encryptString} from 'utils/Common';
import useTestSpecsCrud from './hooks/useTestSpecsCrud';

interface IContext {
  revert: (originalSelector: string) => void;
  add(spec: TTestSpecEntry): void;
  update(selector: string, spec: TTestSpecEntry): void;
  remove(selector: string): void;
  publish(): void;
  runTest(): void;
  cancel(): void;
  dryRun(definitionList: TTestSpecEntry[]): void;
  updateIsInitialized(isInitialized: boolean): void;
  assertionResults?: TAssertionResults;
  specs: TTestSpecEntry[];
  isLoading: boolean;
  isError: boolean;
  isDraftMode: boolean;
  test?: TTest;
  setSelectedSpec(assertionResult?: TAssertionResultEntry): void;
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
  specs: [],
  setSelectedSpec: noop,
  updateIsInitialized: noop,
});

interface IProps {
  testId: string;
  runId: string;
  children: React.ReactNode;
}

export const useTestSpecs = () => useContext(Context);

const TestSpecsProvider = ({children, testId, runId}: IProps) => {
  const dispatch = useAppDispatch();
  const {run} = useTestRun();
  const assertionResults = useAppSelector(state => TestDefinitionSelectors.selectAssertionResults(state));
  const specs = useAppSelector(state => TestDefinitionSelectors.selectDefinitionList(state));
  const isDraftMode = useAppSelector(state => TestDefinitionSelectors.selectIsDraftMode(state));
  const isLoading = useAppSelector(state => TestDefinitionSelectors.selectIsLoading(state));
  const isInitialized = useAppSelector(state => TestDefinitionSelectors.selectIsInitialized(state));
  const {data: test} = useGetTestByIdQuery({testId});

  const {add, cancel, publish, runTest, remove, dryRun, update, init, reset, revert, updateIsInitialized} =
    useTestSpecsCrud({
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
    if (isInitialized && run.state === 'FINISHED') dryRun(specs);
  }, [dryRun, specs, isInitialized, run.state]);

  const setSelectedSpec = useCallback(
    (assertionResult?: TAssertionResultEntry) => {
      dispatch(
        RouterActions.updateSearch({
          [RouterSearchFields.SelectedAssertion]: assertionResult?.selector
            ? encryptString(assertionResult?.selector)
            : '',
        })
      );
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
      specs,
      cancel,
      test,
      setSelectedSpec,
      revert,
      updateIsInitialized,
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
      specs,
      cancel,
      test,
      setSelectedSpec,
      revert,
      updateIsInitialized,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TestSpecsProvider;
