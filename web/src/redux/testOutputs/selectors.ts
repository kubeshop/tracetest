import {createSelector} from '@reduxjs/toolkit';

import {TTestOutput} from 'types/TestOutput.types';
import {RootState} from '../store';
import {endpoints} from '../apis/TraceTest.api';

const testOutputsStateSelector = (state: RootState) => state.testOutputs;

const selectTestRunOutputs = createSelector(
  (state: RootState) => state,
  (state: RootState, testId: string, runId: string) => ({testId, runId}),
  (state, {testId, runId}) => {
    const {data} = endpoints.getRunById.select({testId, runId})(state);

    return data?.outputs ?? [];
  }
);

export const selectTestOutputs = createSelector(
  testOutputsStateSelector,
  selectTestRunOutputs,
  ({outputs}, testRunOutputs) => {
    return outputs.map<TTestOutput>((output, index) => ({
      ...output,
      valueRun: testRunOutputs[index]?.value ?? '',
    }));
  }
);

export const selectTestOutputByIndex = createSelector(
  testOutputsStateSelector,
  (state: RootState, index: number) => index,
  ({outputs}, index) => {
    if (index === -1) return;
    return outputs[index];
  }
);

export const selectIsPending = createSelector(testOutputsStateSelector, ({outputs}) =>
  outputs.some(output => output.isDraft)
);
