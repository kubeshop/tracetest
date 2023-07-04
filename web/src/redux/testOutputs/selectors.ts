import {createSelector} from '@reduxjs/toolkit';
import {RootState} from '../store';

const testOutputsStateSelector = (state: RootState) => state.testOutputs;

export const selectTestOutputs = createSelector(testOutputsStateSelector, ({outputs}) => outputs);

export const selectSelectedOutputs = createSelector(testOutputsStateSelector, ({selectedOutputs}) => selectedOutputs);

export const outputNameSelector = (state: RootState, outputName: string) => outputName;
export const selectIsSelectedOutput = createSelector(
  selectSelectedOutputs,
  outputNameSelector,
  (selectedOutputs, outputName) => !!selectedOutputs.find(selectedOutput => selectedOutput.name === outputName)
);

export const selectIsPending = createSelector(testOutputsStateSelector, ({outputs}) =>
  outputs.some(output => output.isDraft)
);
