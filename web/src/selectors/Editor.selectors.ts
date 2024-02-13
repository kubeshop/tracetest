import {uniqBy} from 'lodash';
import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import AssertionSelectors from './Assertion.selectors';
import SpanSelectors from './Span.selectors';

const stateSelector = (state: RootState) => state;
const paramsSelector = (state: RootState, testId: string, runId: number) => ({
  testId,
  runId,
});

export const selectSelectorAttributeList = createSelector(stateSelector, paramsSelector, (state, {testId, runId}) =>
  AssertionSelectors.selectAllAttributeList(state, testId, runId)
);

export const selectExpressionAttributeList = createSelector(
  stateSelector,
  paramsSelector,
  SpanSelectors.selectMatchedSpans,
  (state, {testId, runId}, spanIds) => {
    const attributeList = AssertionSelectors.selectAttributeList(state, testId, runId, spanIds);

    return uniqBy(attributeList, 'key');
  }
);
