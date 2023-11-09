import {createSlice} from '@reduxjs/toolkit';
import {ICreateTestState, TCreateTestSliceActions} from 'types/Test.types';
import {Plugins} from 'constants/Plugins.constants';
import {SupportedPlugins} from 'constants/Common.constants';

export const initialState: ICreateTestState = {
  draftTest: {
    name: 'Untitled Test',
  },
  stepList: Plugins.REST.stepList,
  stepNumber: 0,
  pluginName: SupportedPlugins.REST,
  isFormValid: true,
};

const createTestSlice = createSlice<ICreateTestState, TCreateTestSliceActions, 'createTest'>({
  name: 'createTest',
  initialState,
  reducers: {
    reset() {
      return initialState;
    },
    setIsFormValid(state, {payload: {isValid}}) {
      state.isFormValid = isValid;
    },
    setPlugin(
      state,
      {
        payload: {
          plugin: {name, stepList},
        },
      }
    ) {
      state.pluginName = name;
      state.stepList = stepList;
      state.draftTest = {};
    },
    setStepNumber(state, {payload: {stepNumber, completeStep = true}}) {
      const currentStep = state.stepList[state.stepNumber];
      if (completeStep) currentStep.status = 'complete';
      else if (currentStep.status !== 'complete') currentStep.status = 'pending';

      const nextStep = state.stepList[stepNumber];
      state.stepNumber = stepNumber;
      if (nextStep && nextStep.status !== 'complete') state.stepList[stepNumber].status = 'selected';
    },
    setDraftTest(state, {payload: {draftTest}}) {
      state.draftTest = {
        ...state.draftTest,
        ...draftTest,
      };
    },
  },
});

export const {setPlugin, setStepNumber, setDraftTest, reset, setIsFormValid} = createTestSlice.actions;
export default createTestSlice.reducer;
