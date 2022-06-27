import {createSlice} from '@reduxjs/toolkit';
import {ICreateTestState, TCreateTestSliceActions} from 'types/Plugins.types';
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
    setPlugin(state, {payload: {plugin}}) {
      state.pluginName = plugin.name;
      state.stepList = plugin.stepList;
      state.draftTest = {};
    },
    setStepNumber(state, {payload: {stepNumber, completeStep = true}}) {
      const currentStep = state.stepList[state.stepNumber];
      if (completeStep) currentStep.status = 'complete';
      else if (currentStep.status !== 'complete') currentStep.status = 'pending';

      const nextStep = state.stepList[stepNumber];
      state.stepNumber = stepNumber;
      if (nextStep.status !== 'complete') state.stepList[stepNumber].status = 'selected';
    },
    setDraftTest(state, {payload: {draftTest}}) {
      state.draftTest = {
        ...state.draftTest,
        ...draftTest,
      };
    },
  },
});

export const {setPlugin, setStepNumber, setDraftTest} = createTestSlice.actions;
export default createTestSlice.reducer;
