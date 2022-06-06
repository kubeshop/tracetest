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
  clearAffectedSpans,
  revertDefinition,
  setSelectedAssertion,
} from '../../../redux/slices/TestDefinition.slice';
import {TAssertionResults} from '../../../types/Assertion.types';
import {TTestDefinitionEntry} from '../../../types/TestDefinition.types';
import useDraftMode from './useDraftMode';
import TestRunGateway from '../../../gateways/TestRun.gateway';

interface IProps {
  runId: string;
  testId: string;
  isDraftMode: boolean;
}

const useTestDefinitionCrud = ({runId, testId, isDraftMode}: IProps) => {
  useDraftMode(isDraftMode);
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const revert = useCallback(
    (originalSelector: string, selector: string) => {
      return dispatch(revertDefinition({originalSelector, selector}));
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
    dispatch(clearAffectedSpans());
    dispatch(setSelectedAssertion(''));

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
