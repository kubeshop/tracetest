import {createSelector} from '@reduxjs/toolkit';

import Span from 'models/Span.model';
import TestRun from 'models/TestRun.model';
import TracetestAPI from 'redux/apis/Tracetest';
import TestRunService from 'services/TestRun.service';
import {RootState} from '../redux/store';

const selectParams = (state: RootState, params: {testId: string; runId: number; spanId: string}) => params;

const selectTestRun = (state: RootState, params: {testId: string; runId: number; spanId: string}) => {
  const {data} = TracetestAPI.instance.endpoints.getRunById.select({testId: params.testId, runId: params.runId})(state);
  return data ?? TestRun({});
};

export const selectSpanById = createSelector([selectTestRun, selectParams], (testRun, params) => {
  const {trace} = testRun;
  return trace?.spans?.find(span => span.id === params.spanId) ?? Span({id: params.spanId});
});

const selectAnalyzerErrors = createSelector([selectTestRun], testRun => {
  const {linter} = testRun;
  return TestRunService.getAnalyzerErrorsHashedBySpan(linter);
});

export const selectAnalyzerErrorsBySpanId = createSelector(
  [selectAnalyzerErrors, selectParams],
  (analyzerErrors, params) => {
    return analyzerErrors[params.spanId];
  }
);

const selectTestSpecs = createSelector([selectTestRun], testRun => {
  const {result} = testRun;
  return TestRunService.getTestSpecsHashedBySpan(result);
});

export const selectTestSpecsBySpanId = createSelector([selectTestSpecs, selectParams], (testSpecs, params) => {
  return testSpecs[params.spanId];
});

const selectTestOutputs = createSelector([selectTestRun], testRun => {
  const {outputs} = testRun;
  return TestRunService.getTestOutputsHashedBySpan(outputs);
});

export const selectTestOutputsBySpanId = createSelector([selectTestOutputs, selectParams], (testOutputs, params) => {
  return testOutputs[params.spanId];
});
