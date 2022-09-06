import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import TestSpecsActions from 'redux/actions/TestSpecs.actions';
import {useAppDispatch} from 'redux/hooks';
import {
  initSpecs,
  resetSpecs,
  addSpec,
  removeSpec,
  revertSpec,
  updateSpec,
  reset as resetAction,
  setIsInitialized,
} from 'redux/slices/TestSpecs.slice';
import {TAssertionResults} from 'types/Assertion.types';
import {TTestSpecEntry} from 'types/TestSpecs.types';
import TestRunGateway from 'gateways/TestRun.gateway';
import useBlockNavigation from 'hooks/useBlockNavigation';
import RouterActions from 'redux/actions/Router.actions';
import {RouterSearchFields} from 'constants/Common.constants';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';

interface IProps {
  runId: string;
  testId: string;
  isDraftMode: boolean;
}

const useTestSpecsCrud = ({runId, testId, isDraftMode}: IProps) => {
  useBlockNavigation(isDraftMode);
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const {onOpen} = useConfirmationModal();

  const revert = useCallback(
    (originalSelector: string) => {
      return dispatch(revertSpec({originalSelector}));
    },
    [dispatch]
  );

  const dryRun = useCallback(
    (definitionList: TTestSpecEntry[]) => {
      return dispatch(TestSpecsActions.dryRun({testId, runId, definitionList})).unwrap();
    },
    [dispatch, runId, testId]
  );

  const publish = useCallback(async () => {
    const {id} = await dispatch(TestSpecsActions.publish({testId, runId})).unwrap();
    dispatch(RouterActions.updateSearch({[RouterSearchFields.SelectedAssertion]: ''}));

    navigate(`/test/${testId}/run/${id}/test`);
  }, [dispatch, navigate, runId, testId]);

  const runTest = useCallback(async () => {
    const {id} = await dispatch(TestRunGateway.runTest(testId)).unwrap();
    dispatch(resetAction());

    navigate(`/test/${testId}/run/${id}`);
  }, [dispatch, navigate, testId]);

  const cancel = useCallback(() => {
    dispatch(resetSpecs());
  }, [dispatch]);

  const add = useCallback(
    async (spec: TTestSpecEntry) => {
      dispatch(addSpec({spec}));
    },
    [dispatch]
  );

  const update = useCallback(
    async (selector: string, spec: TTestSpecEntry) => {
      dispatch(updateSpec({spec, selector}));
    },
    [dispatch]
  );

  const onConfirmRemove = useCallback(
    async (selector: string) => {
      dispatch(removeSpec({selector}));
    },
    [dispatch]
  );
  const remove = useCallback(
    (selector: string) => {
      onOpen('Are you sure you want to remove this test spec?', () => onConfirmRemove(selector));
    },
    [onConfirmRemove, onOpen]
  );

  const init = useCallback(
    (assertionResults: TAssertionResults) => {
      dispatch(initSpecs({assertionResults}));
    },
    [dispatch]
  );

  const updateIsInitialized = useCallback(
    isInitialized => {
      dispatch(setIsInitialized({isInitialized}));
    },
    [dispatch]
  );

  const reset = useCallback(() => {
    dispatch(resetAction());
  }, [dispatch]);

  return {revert, init, updateIsInitialized, reset, add, remove, update, publish, runTest, cancel, dryRun};
};

export default useTestSpecsCrud;
