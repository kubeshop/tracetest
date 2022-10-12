import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import {endpoints} from 'redux/apis/TraceTest.api';
import TestRunService from 'services/TestRun.service';
import {TTestRun} from 'types/TestRun.types';

const stateSelector = (state: RootState) => state;
const paramsSelector = (state: RootState, testId: string, runId: string) => ({testId, runId});

const TestRunSelectors = () => ({
  selectResponseAttributeList: createSelector(stateSelector, paramsSelector, (state, {testId, runId}) => {
    const {data: run = {}} = endpoints.getRunById.select({testId, runId})(state);

    return TestRunService.getResponseAttributeList(run as TTestRun);
  }),
});

export default TestRunSelectors();
