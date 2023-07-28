import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';
import {useCreateTransactionMutation} from 'redux/apis/Tracetest';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {
  initialState,
  reset,
  setDraftTransaction,
  setIsFormValid,
  setStepNumber,
} from 'redux/slices/CreateTransaction.slice';
import CreateTransactionSelectors from 'selectors/CreateTransaction.selectors';
import {ICreateTestStep} from 'types/Plugins.types';
import {ICreateTransactionState, TDraftTransaction} from 'types/Transaction.types';
import useTransactionCrud from '../Transaction/hooks/useTransactionCrud';

interface IContext extends ICreateTransactionState {
  isLoading: boolean;
  stepList: ICreateTestStep[];
  activeStep: string;
  onNext(draft?: TDraftTransaction): void;
  onPrev(): void;
  onCreateTransaction(draft: TDraftTransaction): void;
  onUpdateDraft(draft: TDraftTransaction): void;
  onReset(): void;
  onIsFormValid(isValid: boolean): void;
}

export const Context = createContext<IContext>({
  ...initialState,
  isLoading: false,
  onNext: noop,
  onPrev: noop,
  onCreateTransaction: noop,
  onUpdateDraft: noop,
  onIsFormValid: noop,
  onReset: noop,
  stepList: [],
  activeStep: '',
});

interface IProps {
  children: React.ReactNode;
}

export const useCreateTransaction = () => useContext(Context);

const CreateTransactionProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const [createTransaction, {isLoading: isLoadingCreateTransaction}] = useCreateTransactionMutation();
  const {runTransaction, isEditLoading} = useTransactionCrud();

  const draftTransaction = useAppSelector(CreateTransactionSelectors.selectDraftTransaction);
  const stepNumber = useAppSelector(CreateTransactionSelectors.selectStepNumber);
  const isFormValid = useAppSelector(CreateTransactionSelectors.selectIsFormValid);
  const stepList = useAppSelector(CreateTransactionSelectors.selectStepList);
  const isFinalStep = stepNumber === stepList.length - 1;
  const activeStep = stepList[stepNumber]?.id;

  const onCreateTransaction = useCallback(
    async (draft: TDraftTransaction) => {
      const transaction = await createTransaction(draft).unwrap();
      runTransaction(transaction);
    },
    [createTransaction, runTransaction]
  );

  const onUpdateDraft = useCallback(
    (update: TDraftTransaction) => {
      dispatch(setDraftTransaction({draftTransaction: update}));
    },
    [dispatch]
  );

  const onNext = useCallback(
    (draft: TDraftTransaction = {}) => {
      if (isFinalStep)
        onCreateTransaction({
          ...draftTransaction,
          ...draft,
        });
      else dispatch(setStepNumber({stepNumber: stepNumber + 1}));

      onUpdateDraft(draft);
    },
    [dispatch, draftTransaction, isFinalStep, onCreateTransaction, onUpdateDraft, stepNumber]
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
      draftTransaction,
      stepNumber,
      isLoading: isLoadingCreateTransaction || isEditLoading,
      isFormValid,
      onNext,
      onPrev,
      onCreateTransaction,
      onUpdateDraft,
      onReset,
      stepList,
      activeStep,
      onIsFormValid,
    }),
    [
      draftTransaction,
      stepNumber,
      isLoadingCreateTransaction,
      isEditLoading,
      isFormValid,
      onNext,
      onPrev,
      onCreateTransaction,
      onUpdateDraft,
      onReset,
      stepList,
      activeStep,
      onIsFormValid,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default CreateTransactionProvider;
