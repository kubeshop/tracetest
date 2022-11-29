import {noop} from 'lodash';
import {useNavigate} from 'react-router-dom';
import {outputAdded, outputDeleted, outputsInitiated, outputsReverted, outputUpdated} from 'redux/testOutputs/slice';
import {useAppDispatch, useAppStore} from 'redux/hooks';
import {useReRunMutation, useSetTestOutputsMutation} from 'redux/apis/TraceTest.api';
import {toRawTestOutputs} from 'models/TestOutput.model';
import {createContext, useCallback, useContext, useEffect, useMemo, useState} from 'react';
import {TTestOutput} from 'types/TestOutput.types';
import OutputModal from 'components/OutputModal/OutputModal';
import {selectTestOutputByIndex} from 'redux/testOutputs/selectors';
import {useTest} from '../Test/Test.provider';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';

interface IContext {
  onModalOpen(draft?: TTestOutput): void;
  onModalOpenFromIndex(index: number): void;
  onNavigateAndOpenModal(draft?: TTestOutput): void;
  onCancel(): void;
  onSave(outputs: TTestOutput[]): void;
  onDelete(id: number): void;
  isLoading: boolean;
}

export const Context = createContext<IContext>({
  onModalOpen: noop,
  onModalOpenFromIndex: noop,
  onSave: noop,
  onDelete: noop,
  onCancel: noop,
  onNavigateAndOpenModal: noop,
  isLoading: false,
});

interface IProps {
  children: React.ReactNode;
  testId: string;
  runId: string;
}

export const useTestOutput = () => useContext(Context);

const TestOutputProvider = ({children, testId, runId}: IProps) => {
  const [draft, setDraft] = useState<TTestOutput>();
  const [isOpen, setIsOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const dispatch = useAppDispatch();
  const {getState} = useAppStore();
  const {
    test: {outputs: testOutputs = []},
  } = useTest();
  const [setTestOutputs, {isLoading}] = useSetTestOutputsMutation();
  const [reRunTest, {isLoading: isLoadingReRun}] = useReRunMutation();
  const {onOpen} = useConfirmationModal();
  const navigate = useNavigate();

  useEffect(() => {
    dispatch(outputsInitiated(testOutputs));
  }, [dispatch, testOutputs]);

  const onModalOpen = useCallback((values?: TTestOutput) => {
    setDraft(values);
    setIsOpen(true);
  }, []);

  const onModalOpenFromIndex = useCallback(
    (index: number) => {
      const values = selectTestOutputByIndex(getState(), index);

      setDraft(values!);
      setIsOpen(true);
      setIsEditing(true);
    },
    [getState]
  );

  const onCancel = useCallback(() => {
    dispatch(outputsReverted());
  }, [dispatch]);

  const onDelete = useCallback(
    (index: number) => {
      onOpen(`Are you sure you want to delete the output?`, () => {
        dispatch(outputDeleted(index));
      });
    },
    [dispatch, onOpen]
  );

  const onSubmit = useCallback(
    (values: TTestOutput) => {
      setIsOpen(false);
      if (isEditing) {
        setIsEditing(false);
        dispatch(outputUpdated({output: {...values, id: draft?.id || -1}}));
        return;
      }
      dispatch(outputAdded(values));
    },
    [dispatch, draft?.id, isEditing]
  );

  const onSave = useCallback(
    async (outputs: TTestOutput[]) => {
      await setTestOutputs({testId, testOutputs: toRawTestOutputs(outputs)}).unwrap();
      const run = await reRunTest({testId, runId}).unwrap();

      navigate(`/test/${testId}/run/${run.id}/trigger`);
    },
    [navigate, reRunTest, runId, setTestOutputs, testId]
  );

  const onNavigateAndOpenModal = useCallback(
    async (values?: TTestOutput) => {
      await navigate(`/test/${testId}/run/${runId}/trigger/?tab=outputs`);

      onModalOpen(values);
    },
    [navigate, onModalOpen, runId, testId]
  );

  const value = useMemo<IContext>(
    () => ({
      onModalOpen,
      onModalOpenFromIndex,
      onDelete,
      onSave,
      onCancel,
      isLoading: isLoading || isLoadingReRun,
      onNavigateAndOpenModal,
    }),
    [isLoading, isLoadingReRun, onCancel, onDelete, onModalOpen, onModalOpenFromIndex, onNavigateAndOpenModal, onSave]
  );

  return (
    <Context.Provider value={value}>
      {children}
      <OutputModal
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        onSubmit={onSubmit}
        runId={runId}
        testId={testId}
        output={draft}
        isEditing={isEditing}
      />
    </Context.Provider>
  );
};

export default TestOutputProvider;
