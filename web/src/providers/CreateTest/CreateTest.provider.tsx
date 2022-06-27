import {createContext, useCallback, useContext, useMemo} from 'react';
import {useNavigate} from 'react-router-dom';
import {noop} from 'lodash';
import {ICreateTestState, IPlugin, TDraftTest} from 'types/Plugins.types';
import {initialState, setDraftTest, setPlugin, setStepNumber} from 'redux/slices/CreateTest.slice';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import CreateTestSelectors from 'selectors/CreateTest.selectors';
import {useCreateTestMutation, useRunTestMutation} from 'redux/apis/TraceTest.api';

interface IContext extends ICreateTestState {
  activeStep: string;
  isLoading: boolean;
  onNext(): void;
  onPrev(): void;
  onCreateTest(): void;
  onUpdateDraftTest(draftTest: TDraftTest): void;
  onUpdatePlugin(plugin: IPlugin): void;
  onGoTo(stepId: string): void;
}

export const Context = createContext<IContext>({
  ...initialState,
  activeStep: '',
  isLoading: false,
  onNext: noop,
  onPrev: noop,
  onCreateTest: noop,
  onUpdateDraftTest: noop,
  onUpdatePlugin: noop,
  onGoTo: noop,
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
  const pluginName = useAppSelector(CreateTestSelectors.selectPlugin);
  const activeStep = useAppSelector(CreateTestSelectors.selectActiveStep);

  const onNext = useCallback(() => {
    dispatch(setStepNumber({stepNumber: stepNumber + 1}));
  }, [dispatch, stepNumber]);

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

  const onCreateTest = useCallback(async () => {
    const test = await createTest(draftTest).unwrap();
    const run = await runTest({testId: test.id}).unwrap();

    navigate(`/test/${test.id}/run/${run.id}`);
  }, [createTest, draftTest, navigate, runTest]);

  const onUpdateDraftTest = useCallback(
    (update: TDraftTest) => {
      dispatch(setDraftTest({draftTest: update}));
    },
    [dispatch]
  );

  const onUpdatePlugin = useCallback(
    (plugin: IPlugin) => {
      dispatch(setPlugin({plugin}));
    },
    [dispatch]
  );

  const value = useMemo<IContext>(
    () => ({
      stepList,
      draftTest,
      stepNumber,
      pluginName,
      activeStep,
      isLoading: isLoadingCreateTest || isLoadingRunTest,
      onNext,
      onPrev,
      onGoTo,
      onCreateTest,
      onUpdateDraftTest,
      onUpdatePlugin,
    }),
    [
      stepList,
      draftTest,
      stepNumber,
      pluginName,
      activeStep,
      isLoadingCreateTest,
      isLoadingRunTest,
      onNext,
      onPrev,
      onGoTo,
      onCreateTest,
      onUpdateDraftTest,
      onUpdatePlugin,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default CreateTestProvider;
