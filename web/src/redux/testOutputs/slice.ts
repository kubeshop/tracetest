import {createSlice, PayloadAction} from '@reduxjs/toolkit';

import {TTestOutput} from 'types/TestOutput.types';

interface ITestOutputsState {
  initialOutputs: TTestOutput[];
  outputs: TTestOutput[];
}

const initialState: ITestOutputsState = {
  initialOutputs: [],
  outputs: [],
};

const testOutputsSlice = createSlice({
  name: 'testOutputs',
  initialState,
  reducers: {
    outputsInitiated(state, action: PayloadAction<TTestOutput[]>) {
      state.initialOutputs = action.payload;
      state.outputs = action.payload;
    },
    outputAdded(state, action: PayloadAction<TTestOutput>) {
      state.outputs.push({...action.payload, isDeleted: false, isDraft: true});
    },
    outputUpdated(state, action: PayloadAction<{index: number; output: TTestOutput}>) {
      state.outputs.splice(action.payload.index, 1, {...action.payload.output, isDeleted: false, isDraft: true});
    },
    outputDeleted(state, action: PayloadAction<number>) {
      const output = state.outputs[action.payload];
      if (output) {
        state.outputs.splice(action.payload, 1, {...output, isDeleted: true, isDraft: true});
      }
    },
    outputsReverted(state) {
      state.outputs = state.initialOutputs;
    },
  },
});

export const {outputsInitiated, outputAdded, outputUpdated, outputDeleted, outputsReverted} = testOutputsSlice.actions;

export default testOutputsSlice.reducer;
