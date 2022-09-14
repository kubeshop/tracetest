import {noop} from 'lodash';

import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {createContext, useCallback, useContext, useEffect, useMemo} from 'react';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import {TAssertionResultEntry, TAssertionResults} from 'types/Assertion.types';
import {TTestSpecEntry} from 'types/TestSpecs.types';
import RouterActions from 'redux/actions/Router.actions';
import {RouterSearchFields} from 'constants/Common.constants';
import TestProvider from 'providers/Test';
import useTestSpecsCrud from './hooks/useTestSpecsCrud';

interface IContext {
  revert: (originalSelector: string) => void;
  add(spec: TTestSpecEntry): void;
  update(selector: string, spec: TTestSpecEntry): void;
  remove(selector: string): void;
  publish(): void;
  cancel(): void;
  dryRun(definitionList: TTestSpecEntry[]): void;
  updateIsInitialized(isInitialized: boolean): void;
  assertionResults?: TAssertionResults;
  specs: TTestSpecEntry[];
  isLoading: boolean;
  isError: boolean;
  isDraftMode: boolean;
  setSelectedSpec(assertionResult?: TAssertionResultEntry): void;
}

export const Context = createContext<IContext>({
  add: noop,
  revert: () => noop,
  update: noop,
  remove: noop,
  publish: noop,
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

  const assertionResults = useAppSelector(state => TestSpecsSelectors.selectAssertionResults(state));
  const specs = useAppSelector(state => TestSpecsSelectors.selectSpecs(state));
  const isDraftMode = useAppSelector(state => TestSpecsSelectors.selectIsDraftMode(state));
  const isLoading = useAppSelector(state => TestSpecsSelectors.selectIsLoading(state));
  const isInitialized = useAppSelector(state => TestSpecsSelectors.selectIsInitialized(state));

  const {add, cancel, publish, remove, dryRun, update, init, reset, revert, updateIsInitialized} = useTestSpecsCrud({
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
      const resultList = assertionResults?.resultList || [];
      const positionIndex = resultList.findIndex(({selector}) => selector === assertionResult?.selector);

      dispatch(
        RouterActions.updateSearch({
          [RouterSearchFields.SelectedAssertion]: positionIndex >= 0 ? `${positionIndex}` : '',
        })
      );
    },
    [assertionResults?.resultList, dispatch]
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
      specs,
      cancel,
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
      dryRun,
      assertionResults,
      specs,
      cancel,
      setSelectedSpec,
      revert,
      updateIsInitialized,
    ]
  );

  return (
    <Context.Provider value={value}>
      <TestProvider testId={testId} version={run.testVersion}>
        {children}
      </TestProvider>
    </Context.Provider>
  );
};

export default TestSpecsProvider;
