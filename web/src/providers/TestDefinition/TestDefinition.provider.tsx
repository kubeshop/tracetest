import {noop} from 'lodash';
import {useCallback, createContext, useContext, useMemo} from 'react';
import TestDefinitionActions from '../../redux/actions/TestDefinition.actions';
import {useReRunMutation} from '../../redux/apis/TraceTest.api';
import {useAppDispatch} from '../../redux/hooks';
import {TTestDefinitionEntry} from '../../types/TestDefinition.types';

interface IContext {
  add(testDefinition: TTestDefinitionEntry): void;
  update(selector: string, testDefinition: TTestDefinitionEntry): void;
  remove(selector: string): void;
  isLoading: boolean;
  isError: boolean;
}

export const Context = createContext<IContext>({
  add: noop,
  update: noop,
  remove: noop,
  isLoading: false,
  isError: false,
});

interface IProps {
  testId: string;
  runId: string;
}

export const useTestDefinition = () => useContext(Context);

const TestDefinitionProvider: React.FC<IProps> = ({children, testId, runId}) => {
  const [rerun, {isLoading, isError}] = useReRunMutation();
  const dispatch = useAppDispatch();

  const add = useCallback(
    async (definition: TTestDefinitionEntry) => {
      await dispatch(TestDefinitionActions.add({testId, definition}));
      await rerun({testId, runId});
    },
    [dispatch, rerun, runId, testId]
  );

  const update = useCallback(
    async (selector: string, definition: TTestDefinitionEntry) => {
      await dispatch(TestDefinitionActions.update({testId, definition, selector}));
      await rerun({testId, runId});
    },
    [dispatch, rerun, runId, testId]
  );

  const remove = useCallback(
    async (selector: string) => {
      await dispatch(TestDefinitionActions.remove({testId, selector}));
      await rerun({testId, runId});
    },
    [dispatch, rerun, runId, testId]
  );

  const value = useMemo<IContext>(
    () => ({add, remove, update, isLoading, isError}),
    [add, isError, isLoading, remove, update]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TestDefinitionProvider;
