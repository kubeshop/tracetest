import {createSelector} from '@reduxjs/toolkit';
import {RootState} from '../store';

const testOutputsStateSelector = (state: RootState) => state.testOutputs;

export const selectTestOutputs = createSelector(testOutputsStateSelector, ({outputs}) => outputs);

export const selectSelectedOutputs = createSelector(testOutputsStateSelector, ({selectedOutputs}) => selectedOutputs);

export const outputIdSelector = (state: RootState, outputId: number) => outputId;
export const selectIsSelectedOutput = createSelector(
  selectSelectedOutputs,
  outputIdSelector,
  (selectedOutputs, outputId) => !!selectedOutputs.find(selectedOutput => selectedOutput.id === outputId)
);

export const spanIdSelector = (state: RootState, spanId: string) => spanId;
export const selectOutputsBySpanId = createSelector(selectTestOutputs, spanIdSelector, (outputs, spanId) => {
  return outputs.filter(output => output.spanId === spanId);
});

export const selectIsPending = createSelector(testOutputsStateSelector, ({outputs}) =>
  outputs.some(output => output.isDraft)
);
