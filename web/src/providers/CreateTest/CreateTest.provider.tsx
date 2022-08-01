import {createContext, useCallback, useContext, useMemo} from 'react';
import {useNavigate} from 'react-router-dom';
import {noop} from 'lodash';
import {IPlugin} from 'types/Plugins.types';
import {initialState, setDraftTest, setPlugin, setStepNumber, reset} from 'redux/slices/CreateTest.slice';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import CreateTestSelectors from 'selectors/CreateTest.selectors';
import {useCreateTestMutation, useRunTestMutation} from 'redux/apis/TraceTest.api';
import {ICreateTestState, TDraftTest} from 'types/Test.types';
import TestService from '../../services/Test.service';
import {Plugins} from '../../constants/Plugins.constants';

interface IContext extends ICreateTestState {
  activeStep: string;
  isLoading: boolean;
  plugin: IPlugin;
  onNext(draftTest?: TDraftTest): void;
  onPrev(): void;
  onCreateTest(draftTest: TDraftTest): void;
  onUpdateDraftTest(draftTest: TDraftTest): void;
  onUpdatePlugin(plugin: IPlugin): void;
  onGoTo(stepId: string): void;
  onReset(): void;
}

export const Context = createContext<IContext>({
  ...initialState,
  activeStep: '',
  isLoading: false,
  plugin: Plugins.REST,
  onNext: noop,
  onPrev: noop,
  onCreateTest: noop,
  onUpdateDraftTest: noop,
  onUpdatePlugin: noop,
  onGoTo: noop,
  onReset: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useCreateTest = () => useContext(Context);

const CreateTestProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [createTest, {isLoading: isLoadingCreateTest}] = useCreateTestMutation();
  const [runTest, {isLoading: isLoadingRunTest}] = useRunTestMutation();

  const stepList = useAppSelector(CreateTestSelectors.selectStepList);
  const draftTest = useAppSelector(CreateTestSelectors.selectDraftTest);
  const stepNumber = useAppSelector(CreateTestSelectors.selectStepNumber);
  const plugin = useAppSelector(CreateTestSelectors.selectPlugin);
  const activeStep = useAppSelector(CreateTestSelectors.selectActiveStep);
  const isFinalStep = stepNumber === stepList.length - 1;

  const onCreateTest = useCallback(
    async (draft: TDraftTest) => {
      const rawTest = await TestService.getRequest(plugin, draft);
      const test = await createTest(rawTest).unwrap();
      const run = await runTest({testId: test.id}).unwrap();

      navigate(`/test/${test.id}/run/${run.id}`);
    },
    [createTest, navigate, plugin, runTest]
  );

  const onUpdateDraftTest = useCallback(
    (update: TDraftTest) => {
      dispatch(setDraftTest({draftTest: update}));
    },
    [dispatch]
  );

  const onNext = useCallback(
    (draft: TDraftTest = {}) => {
      if (isFinalStep)
        onCreateTest({
          ...draftTest,
          ...draft,
        });
      else dispatch(setStepNumber({stepNumber: stepNumber + 1}));

      onUpdateDraftTest(draft);
    },
    [dispatch, draftTest, isFinalStep, onCreateTest, onUpdateDraftTest, stepNumber]
  );

  const onPrev = useCallback(() => {
    dispatch(setStepNumber({stepNumber: stepNumber - 1}));
  }, [dispatch, stepNumber]);

  const onGoTo = useCallback(
    (stepId: string) => {
      const stepIndex = stepList.findIndex(({id}) => id === stepId);
      const step = stepList[stepIndex];
      const currentStep = stepList[stepNumber];

      if (step?.status === 'complete' || (currentStep.status === 'complete' && stepIndex === stepNumber + 1))
        dispatch(setStepNumber({stepNumber: stepIndex, completeStep: false}));
    },
    [dispatch, stepList, stepNumber]
  );

  const onUpdatePlugin = useCallback(
    (newPlugin: IPlugin) => {
      dispatch(setPlugin({plugin: newPlugin}));
    },
    [dispatch]
  );

  const onReset = useCallback(() => {
    dispatch(reset());
  }, [dispatch]);

  const value = useMemo<IContext>(
    () => ({
      stepList,
      draftTest,
      stepNumber,
      pluginName: plugin.name,
      plugin,
      activeStep,
      isLoading: isLoadingCreateTest || isLoadingRunTest,
      onNext,
      onPrev,
      onGoTo,
      onCreateTest,
      onUpdateDraftTest,
      onUpdatePlugin,
      onReset,
    }),
    [
      stepList,
      draftTest,
      stepNumber,
      plugin,
      activeStep,
      isLoadingCreateTest,
      isLoadingRunTest,
      onNext,
      onPrev,
      onGoTo,
      onCreateTest,
      onUpdateDraftTest,
      onUpdatePlugin,
      onReset,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default CreateTestProvider;
