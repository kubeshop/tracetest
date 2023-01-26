import {createSlice, PayloadAction} from '@reduxjs/toolkit';

import {TTestOutput, TTestRunOutput} from 'types/TestOutput.types';

interface ITestOutputsState {
  initialOutputs: TTestOutput[];
  outputs: TTestOutput[];
  selectedOutputs: TTestOutput[];
}

const initialState: ITestOutputsState = {
  initialOutputs: [],
  outputs: [],
  selectedOutputs: [],
};

const testOutputsSlice = createSlice({
  name: 'testOutputs',
  initialState,
  reducers: {
    outputsReseted() {
      return initialState;
    },
    outputsInitiated(state, {payload: outputs}: PayloadAction<TTestOutput[]>) {
      state.initialOutputs = outputs;
      state.outputs = outputs;
    },
    outputAdded(state, {payload: outputs}: PayloadAction<TTestOutput>) {
      state.outputs.push({...outputs, isDeleted: false, isDraft: true, id: state.outputs.length});
    },
    outputUpdated(state, {payload: {output}}: PayloadAction<{output: TTestOutput}>) {
      state.outputs.splice(output.id, 1, {...output, isDeleted: false, isDraft: true});
    },
    outputDeleted(state, {payload: outputId}: PayloadAction<number>) {
      const output = state.outputs[outputId];
      if (output) {
        state.outputs.splice(outputId, 1, {...output, isDeleted: true, isDraft: true});
      }
    },
    outputsReverted(state) {
      state.outputs = state.initialOutputs;
    },
    outputsSelectedOutputsChanged(state, {payload: outputs}: PayloadAction<TTestOutput[]>) {
      state.selectedOutputs = outputs;
    },
    outputsTestRunOutputsMerged(state, {payload: runOutputs}: PayloadAction<TTestRunOutput[]>) {
      state.outputs = state.outputs.map((output, index) => ({
        ...output,
        valueRun: runOutputs[index]?.value ?? '',
        spanId: runOutputs[index]?.spanId ?? '',
      }));
    },
  },
});

export const {
  outputsInitiated,
  outputAdded,
  outputUpdated,
  outputDeleted,
  outputsReverted,
  outputsReseted,
  outputsSelectedOutputsChanged,
  outputsTestRunOutputsMerged,
} = testOutputsSlice.actions;

export default testOutputsSlice.reducer;
