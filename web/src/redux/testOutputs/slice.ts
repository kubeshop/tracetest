import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import TestOutput from 'models/TestOutput.model';
import TestRunOutput from 'models/TestRunOutput.model';

interface ITestOutputsState {
  initialOutputs: TestOutput[];
  outputs: TestOutput[];
  selectedOutputs: TestOutput[];
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
    outputsInitiated(state, {payload: outputs}: PayloadAction<TestOutput[]>) {
      state.initialOutputs = outputs;
      state.outputs = outputs;
    },
    outputAdded(state, {payload: outputs}: PayloadAction<TestOutput>) {
      state.outputs.push({...outputs, isDeleted: false, isDraft: true, id: state.outputs.length});
    },
    outputUpdated(state, {payload: {output}}: PayloadAction<{output: TestOutput}>) {
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
    outputsSelectedOutputsChanged(state, {payload: outputs}: PayloadAction<TestOutput[]>) {
      state.selectedOutputs = outputs;
    },
    outputsTestRunOutputsMerged(state, {payload: runOutputs}: PayloadAction<TestRunOutput[]>) {
      state.outputs = state.outputs.map((output, index) => ({
        ...output,
        valueRun: runOutputs[index]?.value ?? '',
        spanId: runOutputs[index]?.spanId ?? '',
        error: runOutputs[index]?.error ?? '',
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
