import {createSlice} from '@reduxjs/toolkit';
import {ICreateTestState, TCreateTestSliceActions} from 'types/Test.types';
import {Plugins, SupportedPlugins} from 'constants/Plugins.constants';

export const initialState: ICreateTestState = {
  draftTest: {},
  stepList: Plugins.REST.stepList,
  stepNumber: 0,
  pluginName: SupportedPlugins.REST,
};

const createTestSlice = createSlice<ICreateTestState, TCreateTestSliceActions, 'createTest'>({
  name: 'createTest',
  initialState,
  reducers: {
    reset() {
      return initialState;
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

export const {setPlugin, setStepNumber, setDraftTest, reset} = createTestSlice.actions;
export default createTestSlice.reducer;
