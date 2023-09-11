import {noop} from 'lodash';
import {createContext, useCallback, useContext, useEffect, useMemo, useState} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {selectIsPending, selectSelectedOutputs, selectTestOutputs} from 'redux/testOutputs/selectors';
import {
  outputAdded,
  outputDeleted,
  outputsInitiated,
  outputsReseted,
  outputsReverted,
  outputsSelectedOutputsChanged,
  outputsTestRunOutputsMerged,
  outputUpdated,
} from 'redux/testOutputs/slice';
import TestOutput from 'models/TestOutput.model';
import SpanSelectors from 'selectors/Span.selectors';
import useValidateOutput from './hooks/useValidateOutput';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';
import {useDashboard} from '../Dashboard/Dashboard.provider';
import {useVariableSet} from '../VariableSet';
import {useTest} from '../Test/Test.provider';
import {useTestRun} from '../TestRun/TestRun.provider';

const {useParseExpressionMutation} = TracetestAPI.instance;

interface IContext {
  isDraftMode: boolean;
  isEditing: boolean;
  isLoading: boolean;
  isOpen: boolean;
  isValid: boolean;
  onCancel(): void;
  onClose(): void;
  onDelete(id: number): void;
  onNavigateAndOpen(draft?: TestOutput): void;
  onOpen(draft?: TestOutput): void;
  onSubmit(values: TestOutput): void;
  onSelectedOutputs(outputs: TestOutput[]): void;
  onValidate(_: any, output: TestOutput): void;
  output?: TestOutput;
  outputs: TestOutput[];
  selectedOutputs: TestOutput[];
}

export const Context = createContext<IContext>({
  isDraftMode: false,
  isEditing: false,
  isLoading: false,
  isOpen: false,
  isValid: false,
  onCancel: noop,
  onClose: noop,
  onDelete: noop,
  onNavigateAndOpen: noop,
  onOpen: noop,
  onSubmit: noop,
  onSelectedOutputs: noop,
  onValidate: noop,
  output: undefined,
  outputs: [],
  selectedOutputs: [],
});

interface IProps {
  children: React.ReactNode;
  runId: number;
  testId: string;
}

export const useTestOutput = () => useContext(Context);

const TestOutputProvider = ({children, runId, testId}: IProps) => {
  const dispatch = useAppDispatch();
  const [parseExpressionMutation, {isLoading}] = useParseExpressionMutation();
  const [draft, setDraft] = useState<TestOutput>();
  const [isEditing, setIsEditing] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const {onOpen: onOpenConfirmationModal} = useConfirmationModal();
  const {navigate} = useDashboard();
  const {selectedVariableSet} = useVariableSet();
  const {
    test: {outputs: testOutputs = []},
  } = useTest();
  const {
    run: {outputs: runOutputs = []},
  } = useTestRun();
  const outputs = useAppSelector(state => selectTestOutputs(state));
  const selectedOutputs = useAppSelector(selectSelectedOutputs);
  const isDraftMode = useAppSelector(selectIsPending);
  const spanIdList = useAppSelector(SpanSelectors.selectMatchedSpans);
  const {isValid, onValidate} = useValidateOutput({spanIdList});

  useEffect(() => {
    dispatch(outputsInitiated(testOutputs));

    return () => {
      dispatch(outputsReseted());
    };
  }, [dispatch, testOutputs]);

  useEffect(() => {
    dispatch(outputsTestRunOutputsMerged(runOutputs));
  }, [dispatch, runOutputs]);

  const handleOpen = useCallback((values?: TestOutput) => {
    setDraft(values);
    setIsOpen(true);
    const id = values?.id ?? -1;
    setIsEditing(id !== -1);
  }, []);

  const onOpen = useCallback(
    (values?: TestOutput) => {
      if (isValid) {
        onOpenConfirmationModal({
          title: 'Unsaved changes',
          heading: 'Discard unsaved changes?',
          okText: 'Discard',
          onConfirm: () => {
            handleOpen(values);
          },
        });
      } else handleOpen(values);
    },
    [handleOpen, isValid, onOpenConfirmationModal]
  );

  const onClose = useCallback(() => {
    setDraft(undefined);
    setIsOpen(false);
    onValidate(undefined, TestOutput({}));
  }, [onValidate]);

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
    async (values: TestOutput, matchedSpanId?: string) => {
      const spanId = values.spanId || matchedSpanId || '';
      const props = {
        expression: values.value,
        context: {
          testId,
          runId,
          spanId,
          selector: values.selector,
          variableSetId: selectedVariableSet?.id,
        },
      };
      const parsedExpression = await parseExpressionMutation(props).unwrap();
      const valueRunDraft = parsedExpression?.[0] ?? '';

      setIsOpen(false);
      if (isEditing) {
        setIsEditing(false);
        dispatch(outputUpdated({output: {...values, spanId, valueRunDraft, id: draft?.id ?? -1}}));
        return;
      }
      dispatch(outputAdded({...values, valueRunDraft, spanId}));
    },
    [dispatch, draft?.id, isEditing, parseExpressionMutation, runId, selectedVariableSet?.id, testId]
  );

  const onNavigateAndOpen = useCallback(
    async (values?: TestOutput) => {
      await navigate(`/test/${testId}/run/${runId}/test/?tab=outputs`);
      onOpen(values);
    },
    [navigate, onOpen, runId, testId]
  );

  const onSelectedOutputs = useCallback(
    (outputList: TestOutput[]) => {
      navigate(`/test/${testId}/run/${runId}/test/?tab=outputs`);
      dispatch(outputsSelectedOutputsChanged(outputList));
    },
    [dispatch, navigate, runId, testId]
  );

  const value = useMemo<IContext>(
    () => ({
      isDraftMode,
      isEditing,
      isLoading,
      isOpen,
      isValid,
      onValidate,
      onCancel,
      onClose,
      onDelete,
      onNavigateAndOpen,
      onOpen,
      onSubmit,
      onSelectedOutputs,
      output: draft,
      outputs,
      selectedOutputs,
    }),
    [
      draft,
      isDraftMode,
      isEditing,
      isLoading,
      isOpen,
      isValid,
      onCancel,
      onClose,
      onDelete,
      onNavigateAndOpen,
      onOpen,
      onSelectedOutputs,
      onSubmit,
      onValidate,
      outputs,
      selectedOutputs,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TestOutputProvider;
