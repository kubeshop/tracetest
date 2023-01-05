import {noop} from 'lodash';
import {createContext, useCallback, useContext, useEffect, useMemo, useState} from 'react';
import {useNavigate} from 'react-router-dom';

import {useParseExpressionMutation} from 'redux/apis/TraceTest.api';
import {useAppDispatch} from 'redux/hooks';
import {outputAdded, outputDeleted, outputsInitiated, outputsReverted, outputUpdated} from 'redux/testOutputs/slice';
import {TTestOutput} from 'types/TestOutput.types';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';
import {useEnvironment} from '../Environment/Environment.provider';
import {useTest} from '../Test/Test.provider';

interface IContext {
  isEditing: boolean;
  isLoading: boolean;
  isOpen: boolean;
  onCancel(): void;
  onClose(): void;
  onDelete(id: number): void;
  onNavigateAndOpen(draft?: TTestOutput): void;
  onOpen(draft?: TTestOutput): void;
  onSubmit(values: TTestOutput): void;
  output?: TTestOutput;
}

export const Context = createContext<IContext>({
  isEditing: false,
  isLoading: false,
  isOpen: false,
  onCancel: noop,
  onClose: noop,
  onDelete: noop,
  onNavigateAndOpen: noop,
  onOpen: noop,
  onSubmit: noop,
  output: undefined,
});

interface IProps {
  children: React.ReactNode;
  runId: string;
  testId: string;
}

export const useTestOutput = () => useContext(Context);

const TestOutputProvider = ({children, runId, testId}: IProps) => {
  const dispatch = useAppDispatch();
  const [parseExpressionMutation, {isLoading}] = useParseExpressionMutation();
  const [draft, setDraft] = useState<TTestOutput>();
  const [isEditing, setIsEditing] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const {onOpen: onOpenConfirmationModal} = useConfirmationModal();
  const navigate = useNavigate();
  const {selectedEnvironment} = useEnvironment();
  const {
    test: {outputs: testOutputs = []},
  } = useTest();

  useEffect(() => {
    dispatch(outputsInitiated(testOutputs));
  }, [dispatch, testOutputs]);

  const onOpen = useCallback((values?: TTestOutput) => {
    setDraft(values);
    setIsOpen(true);
    const id = values?.id ?? -1;
    setIsEditing(id !== -1);
  }, []);

  const onClose = useCallback(() => {
    setDraft(undefined);
    setIsOpen(false);
  }, []);

  const onCancel = useCallback(() => {
    dispatch(outputsReverted());
  }, [dispatch]);

  const onDelete = useCallback(
    (index: number) => {
      onOpenConfirmationModal({
        title: `Are you sure you want to delete the output?`,
        onConfirm: () => {
          dispatch(outputDeleted(index));
        },
      });
    },
    [dispatch, onOpenConfirmationModal]
  );

  const onSubmit = useCallback(
    async (values: TTestOutput, spanId?: string) => {
      const props = {
        expression: values.value,
        context: {
          testId,
          runId,
          spanId: spanId ?? '',
          selector: values.selector,
          environmentId: selectedEnvironment?.id,
        },
      };
      const parsedExpression = await parseExpressionMutation(props).unwrap();
      const valueRunDraft = parsedExpression?.[0] ?? '';

      setIsOpen(false);
      if (isEditing) {
        setIsEditing(false);
        dispatch(outputUpdated({output: {...values, valueRunDraft, id: draft?.id ?? -1}}));
        return;
      }
      dispatch(outputAdded({...values, valueRunDraft}));
    },
    [dispatch, draft?.id, isEditing, parseExpressionMutation, runId, selectedEnvironment?.id, testId]
  );

  const onNavigateAndOpen = useCallback(
    async (values?: TTestOutput) => {
      await navigate(`/test/${testId}/run/${runId}/test/?tab=outputs`);
      onOpen(values);
    },
    [navigate, onOpen, runId, testId]
  );

  const value = useMemo<IContext>(
    () => ({
      isEditing,
      isLoading,
      isOpen,
      onCancel,
      onClose,
      onDelete,
      onNavigateAndOpen,
      onOpen,
      onSubmit,
      output: draft,
    }),
    [draft, isEditing, isLoading, isOpen, onCancel, onClose, onDelete, onNavigateAndOpen, onOpen, onSubmit]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TestOutputProvider;
