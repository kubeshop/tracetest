import {createSlice} from '@reduxjs/toolkit';
import {ICreateTransactionState, TCreateTransactionSliceActions} from 'types/Transaction.types';

export const initialState: ICreateTransactionState = {
  draftTransaction: {},
  stepNumber: 0,
  isFormValid: false,
  stepList: [
    {
      id: 'basic-details',
      name: 'General Info',
      title: 'General Info',
      component: 'BasicDetails',
    },
    {
      id: 'tests-selection',
      name: 'Tests',
      title: 'Tests',
      component: 'TestsSelection',
      isDefaultValid: true,
    },
  ],
};

const createTransactionSlice = createSlice<
  ICreateTransactionState,
  TCreateTransactionSliceActions,
  'createTransaction'
>({
  name: 'createTransaction',
  initialState,
  reducers: {
    reset() {
      return initialState;
    },
    setIsFormValid(state, {payload: {isValid}}) {
      state.isFormValid = isValid;
    },
    setStepNumber(state, {payload: {stepNumber, completeStep = true}}) {
      const currentStep = state.stepList[state.stepNumber];
      if (completeStep) currentStep.status = 'complete';
      else if (currentStep.status !== 'complete') currentStep.status = 'pending';

      const nextStep = state.stepList[stepNumber];
      state.stepNumber = stepNumber;
      if (nextStep && nextStep.status !== 'complete') state.stepList[stepNumber].status = 'selected';
    },
    setDraftTransaction(state, {payload: {draftTransaction}}) {
      state.draftTransaction = {
        ...state.draftTransaction,
        ...draftTransaction,
      };
    },
  },
});

export const {setStepNumber, setDraftTransaction, reset, setIsFormValid} = createTransactionSlice.actions;
export default createTransactionSlice.reducer;
