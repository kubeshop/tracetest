import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import TestDefinitionActions from '../../../redux/actions/TestDefinition.actions';
import {useAppDispatch} from '../../../redux/hooks';
import {
  addDefinition,
  initDefinitionList,
  removeDefinition,
  resetDefinitionList,
  updateDefinition,
  reset as resetAction,
  revertDefinition,
  setSelectedAssertion,
} from '../../../redux/slices/TestDefinition.slice';
import {TAssertionResults} from '../../../types/Assertion.types';
import {TTestDefinitionEntry} from '../../../types/TestDefinition.types';
import TestRunGateway from '../../../gateways/TestRun.gateway';
import useBlockNavigation from '../../../hooks/useBlockNavigation';

interface IProps {
  runId: string;
  testId: string;
  isDraftMode: boolean;
}

const useTestDefinitionCrud = ({runId, testId, isDraftMode}: IProps) => {
  useBlockNavigation(isDraftMode);
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const revert = useCallback(
    (originalSelector: string) => {
      return dispatch(revertDefinition({originalSelector}));
    },
    [dispatch]
  );

  const dryRun = useCallback(
    (definitionList: TTestDefinitionEntry[]) => {
      return dispatch(TestDefinitionActions.dryRun({testId, runId, definitionList})).unwrap();
    },
    [dispatch, runId, testId]
  );

  const publish = useCallback(async () => {
    const {id} = await dispatch(TestDefinitionActions.publish({testId, runId})).unwrap();
    dispatch(setSelectedAssertion());

    navigate(`/test/${testId}/run/${id}`);
  }, [dispatch, navigate, runId, testId]);

  const runTest = useCallback(async () => {
    const {id} = await dispatch(TestRunGateway.runTest(testId)).unwrap();
    dispatch(resetAction());

    navigate(`/test/${testId}/run/${id}`);
  }, [dispatch, navigate, testId]);

  const cancel = useCallback(() => {
    dispatch(resetDefinitionList());
  }, [dispatch]);

  const add = useCallback(
    async (definition: TTestDefinitionEntry) => {
      dispatch(addDefinition({definition}));
    },
    [dispatch]
  );

  const update = useCallback(
    async (selector: string, definition: TTestDefinitionEntry) => {
      dispatch(updateDefinition({definition, selector}));
    },
    [dispatch]
  );

  const remove = useCallback(
    async (selector: string) => {
      dispatch(removeDefinition({selector}));
    },
    [dispatch]
  );

  const init = useCallback(
    (assertionResults: TAssertionResults) => {
      dispatch(initDefinitionList({assertionResults}));
    },
    [dispatch]
  );

  const reset = useCallback(() => {
    dispatch(resetAction());
  }, [dispatch]);

  return {revert, init, reset, add, remove, update, publish, runTest, cancel, dryRun};
};

export default useTestDefinitionCrud;
