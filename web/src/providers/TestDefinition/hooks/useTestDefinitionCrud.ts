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
} from '../../../redux/slices/TestDefinition.slice';
import {TAssertionResults} from '../../../types/Assertion.types';
import {TTestDefinitionEntry} from '../../../types/TestDefinition.types';
import useDraftMode from './useDraftMode';

interface IProps {
  runId: string;
  testId: string;
}

const useTestDefinitionCrud = ({runId, testId}: IProps) => {
  const {isDraftMode, setIsDraftMode} = useDraftMode();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const dryRun = useCallback(
    (definitionList: TTestDefinitionEntry[]) => {
      return dispatch(TestDefinitionActions.dryRun({testId, runId, definitionList})).unwrap();
    },
    [dispatch, runId, testId]
  );

  const publish = useCallback(async () => {
    const {id} = await dispatch(TestDefinitionActions.publish({testId, runId})).unwrap();
    setIsDraftMode(false);
    dispatch(clearAffectedSpans());

    navigate(`/test/${testId}/run/${id}`);
  }, [dispatch, navigate, runId, setIsDraftMode, testId]);

  const cancel = useCallback(() => {
    setIsDraftMode(false);
    dispatch(resetDefinitionList());
  }, [dispatch, setIsDraftMode]);

  const add = useCallback(
    async (definition: TTestDefinitionEntry) => {
      dispatch(addDefinition({definition}));
      setIsDraftMode(true);
    },
    [dispatch, setIsDraftMode]
  );

  const update = useCallback(
    async (selector: string, definition: TTestDefinitionEntry) => {
      dispatch(updateDefinition({definition, selector}));
      setIsDraftMode(true);
    },
    [dispatch, setIsDraftMode]
  );

  const remove = useCallback(
    async (selector: string) => {
      dispatch(removeDefinition({selector}));
      setIsDraftMode(true);
    },
    [dispatch, setIsDraftMode]
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

  return {init, reset, add, remove, update, publish, cancel, dryRun, isDraftMode};
};

export default useTestDefinitionCrud;
