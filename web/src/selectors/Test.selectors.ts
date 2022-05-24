import {createSelector} from '@reduxjs/toolkit';

import {endpoints} from 'redux/apis/TraceTest.api';
import {RootState} from '../redux/store';

const stateSelector = (state: RootState) => state;

const testIdSelector = (state: RootState, testId: string) => testId;

const TestSelectors = () => ({
  selectTest: createSelector(stateSelector, testIdSelector, (state, testId) => {
    const {data: test} = endpoints.getTestById.select({testId})(state);

    return test!;
  }),
});

export default TestSelectors();
