import {createSlice} from '@reduxjs/toolkit';
import {ICreateTestSuiteState, TCreateTestSuiteSliceActions} from 'types/TestSuite.types';

export const initialState: ICreateTestSuiteState = {
  draft: {},
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

const createTestSuiteSlice = createSlice<
  ICreateTestSuiteState,
  TCreateTestSuiteSliceActions,
  'createTestSuite'
>({
  name: 'createTestSuite',
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
    setDraft(state, {payload: {draft}}) {
      state.draft = {
        ...state.draft,
        ...draft,
      };
    },
  },
});

export const {setStepNumber, setDraft, reset, setIsFormValid} = createTestSuiteSlice.actions;
export default createTestSuiteSlice.reducer;
