import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';

import AssertionResults from 'models/AssertionResults.model';
import Test from 'models/Test.model';
import TestSpecs, {TTestSpecEntry} from 'models/TestSpecs.model';
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
  setSelectorSuggestions as setSelectorSuggestionsAction,
  setPrevSelector as setPrevSelectorAction,
} from 'redux/slices/TestSpecs.slice';
import {ISuggestion} from 'types/TestSpecs.types';
import useBlockNavigation from 'hooks/useBlockNavigation';
import RouterActions from 'redux/actions/Router.actions';
import {RouterSearchFields} from 'constants/Common.constants';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {useNotification} from 'providers/Notification/Notification.provider';

interface IProps {
  runId: string;
  testId: string;
  test: Test;
  isDraftMode: boolean;
  assertionResults?: AssertionResults;
}

const useTestSpecsCrud = ({runId, testId, test, isDraftMode, assertionResults}: IProps) => {
  useBlockNavigation(isDraftMode);
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const {onOpen} = useConfirmationModal();
  const {showNotification} = useNotification();

  const setSelectedSpec = useCallback(
    selector => {
      const resultList = assertionResults?.resultList || [];
      const positionIndex = resultList.findIndex(result => result.selector === selector);

      dispatch(
        RouterActions.updateSearch({
          [RouterSearchFields.SelectedAssertion]: positionIndex >= 0 ? `${positionIndex}` : undefined,
        })
      );
    },
    [assertionResults?.resultList, dispatch]
  );

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
    const {id} = await dispatch(TestSpecsActions.publish({test, testId, runId})).unwrap();
    dispatch(RouterActions.updateSearch({[RouterSearchFields.SelectedAssertion]: ''}));

    showNotification({
      type: 'success',
      title: 'Your test has been saved successfully',
      description: 'A new test run has been generated.',
    });

    navigate(`/test/${testId}/run/${id}/test`);
  }, [dispatch, navigate, runId, showNotification, test, testId]);

  const cancel = useCallback(() => {
    dispatch(resetSpecs());
  }, [dispatch]);

  const add = useCallback(
    async (spec: TTestSpecEntry) => {
      dispatch(addSpec({spec}));
      setSelectedSpec(spec.selector);
    },
    [dispatch, setSelectedSpec]
  );

  const update = useCallback(
    async (selector: string, spec: TTestSpecEntry) => {
      dispatch(updateSpec({spec, selector}));

      setSelectedSpec(selector);
    },
    [dispatch, setSelectedSpec]
  );

  const onConfirmRemove = useCallback(
    async (selector: string) => {
      dispatch(removeSpec({selector}));
    },
    [dispatch]
  );
  const remove = useCallback(
    (selector: string) => {
      onOpen({
        title: 'Are you sure you want to remove this test spec?',
        onConfirm: () => onConfirmRemove(selector),
      });
    },
    [onConfirmRemove, onOpen]
  );

  const init = useCallback(
    (initialAssertionResults: AssertionResults, specs: TestSpecs) => {
      dispatch(initSpecs({assertionResults: initialAssertionResults, specs}));
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

  const setSelectorSuggestions = useCallback(
    (selectorSuggestions: ISuggestion[]) => {
      dispatch(setSelectorSuggestionsAction(selectorSuggestions));
    },
    [dispatch]
  );

  const setPrevSelector = useCallback(
    (prevSelector: string) => {
      dispatch(setPrevSelectorAction({prevSelector}));
    },
    [dispatch]
  );

  return {
    revert,
    init,
    updateIsInitialized,
    reset,
    add,
    remove,
    update,
    publish,
    cancel,
    dryRun,
    setSelectedSpec,
    setSelectorSuggestions,
    setPrevSelector,
  };
};

export default useTestSpecsCrud;
