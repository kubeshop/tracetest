import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {initialState, reset, setDraft, setIsFormValid, setStepNumber} from 'redux/slices/CreateTestSuite.slice';
import CreateTestSuitesSelectors from 'selectors/CreateTestSuite.selectors';
import {ICreateTestStep} from 'types/Plugins.types';
import {ICreateTestSuiteState, TDraftTestSuite} from 'types/TestSuite.types';
import {useTestSuiteCrud} from '../TestSuite';

interface IContext extends ICreateTestSuiteState {
  isLoading: boolean;
  stepList: ICreateTestStep[];
  activeStep: string;
  onNext(draft?: TDraftTestSuite): void;
  onPrev(): void;
  onCreate(draft: TDraftTestSuite): void;
  onUpdate(draft: TDraftTestSuite): void;
  onReset(): void;
  onIsFormValid(isValid: boolean): void;
}

export const Context = createContext<IContext>({
  ...initialState,
  isLoading: false,
  onNext: noop,
  onPrev: noop,
  onCreate: noop,
  onUpdate: noop,
  onIsFormValid: noop,
  onReset: noop,
  stepList: [],
  activeStep: '',
});

interface IProps {
  children: React.ReactNode;
}

export const useCreateTestSuite = () => useContext(Context);

const CreateTestSuiteProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const {isEditLoading, create, isLoadingCreate} = useTestSuiteCrud();

  const draft = useAppSelector(CreateTestSuitesSelectors.selectDraft);
  const stepNumber = useAppSelector(CreateTestSuitesSelectors.selectStepNumber);
  const isFormValid = useAppSelector(CreateTestSuitesSelectors.selectIsFormValid);
  const stepList = useAppSelector(CreateTestSuitesSelectors.selectStepList);
  const isFinalStep = stepNumber === stepList.length - 1;
  const activeStep = stepList[stepNumber]?.id;

  const onCreate = useCallback(
    async (values: TDraftTestSuite) => {
      await create(values);
    },
    [create]
  );

  const onUpdate = useCallback(
    (update: TDraftTestSuite) => {
      dispatch(setDraft({draft: update}));
    },
    [dispatch]
  );

  const onNext = useCallback(
    (values: TDraftTestSuite = {}) => {
      if (isFinalStep)
        onCreate({
          ...draft,
          ...values,
        });
      else dispatch(setStepNumber({stepNumber: stepNumber + 1}));

      onUpdate(values);
    },
    [isFinalStep, onCreate, draft, dispatch, stepNumber, onUpdate]
  );

  const onPrev = useCallback(() => {
    dispatch(setStepNumber({stepNumber: stepNumber - 1}));
  }, [dispatch, stepNumber]);

  const onReset = useCallback(() => {
    dispatch(reset());
  }, [dispatch]);

  const onIsFormValid = useCallback(
    (isValid: boolean) => {
      dispatch(setIsFormValid({isValid}));
    },
    [dispatch]
  );

  const value = useMemo<IContext>(
    () => ({
      draft,
      stepNumber,
      isLoading: isLoadingCreate || isEditLoading,
      isFormValid,
      onNext,
      onPrev,
      onCreate,
      onUpdate,
      onReset,
      stepList,
      activeStep,
      onIsFormValid,
    }),
    [
      draft,
      stepNumber,
      isLoadingCreate,
      isEditLoading,
      isFormValid,
      onNext,
      onPrev,
      onCreate,
      onUpdate,
      onReset,
      stepList,
      activeStep,
      onIsFormValid,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default CreateTestSuiteProvider;
