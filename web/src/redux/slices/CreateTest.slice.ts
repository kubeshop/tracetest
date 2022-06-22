import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {ICreateTestState, IPlugin, TDraftTest} from 'types/Plugins.types';
import {Plugins, SupportedPlugins} from 'constants/Plugins.constants';

export const initialState: ICreateTestState = {
  draftTest: {},
  stepList: Plugins.REST.stepList,
  stepNumber: 0,
  pluginName: SupportedPlugins.REST,
};

const createTestSlice = createSlice({
  name: 'createTest',
  initialState,
  reducers: {
    setPlugin(state, {payload: {plugin}}: PayloadAction<{plugin: IPlugin}>) {
      state.pluginName = plugin.name;
      state.stepList = plugin.stepList;
      state.draftTest = {};
    },
    setStepNumber(
      state,
      {payload: {stepNumber, completeStep = true}}: PayloadAction<{stepNumber: number; completeStep?: boolean}>
    ) {
      const currentStep = state.stepList[state.stepNumber];
      if (completeStep) currentStep.status = 'complete';
      else if (currentStep.status !== 'complete') currentStep.status = 'pending';

      const nextStep = state.stepList[stepNumber];
      state.stepNumber = stepNumber;
      if (nextStep.status !== 'complete') state.stepList[stepNumber].status = 'selected';
    },
    setDraftTest(state, {payload: {draftTest}}: PayloadAction<{draftTest: TDraftTest}>) {
      state.draftTest = {
        ...state.draftTest,
        ...draftTest,
      };
    },
  },
});

export const {setPlugin, setStepNumber, setDraftTest} = createTestSlice.actions;
export default createTestSlice.reducer;
