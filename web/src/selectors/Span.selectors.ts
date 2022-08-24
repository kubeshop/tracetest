import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import {endpoints} from '../redux/apis/TraceTest.api';

const spansStateSelector = (state: RootState) => state.spans;
const stateSelector = (state: RootState) => state;
const paramsSelector = (state: RootState, spanId: string, testId: string, runId: string) => ({spanId, testId, runId});

const selectAffectedSpans = createSelector(
  spansStateSelector,
  stateSelector,
  ({affectedSpans}, {testDefinition: {assertionResults, selectedAssertion}}) => {
    if (!selectedAssertion) return affectedSpans;

    const foundAssertion = assertionResults?.resultList.find(({selector}) => selector === selectedAssertion);

    return !foundAssertion ? [] : affectedSpans;
  }
);
const SpanSelectors = () => ({
  selectAffectedSpans,
  selectSpanById: createSelector(stateSelector, paramsSelector, (state, {spanId, testId, runId}) => {
    const {data: {trace} = {}} = endpoints.getRunById.select({testId, runId})(state);

    const spanList = trace?.spans || [];

    return spanList.find(span => span.id === spanId);
  }),
  selectSelectedSpan: createSelector(spansStateSelector, ({selectedSpan}) => selectedSpan),
  selectFocusedSpan: createSelector(spansStateSelector, ({focusedSpan}) => focusedSpan),
  selectMatchedSpans: createSelector(spansStateSelector, ({matchedSpans}) => matchedSpans),
  selectSearchText: createSelector(spansStateSelector, ({searchText}) => searchText),
});

export default SpanSelectors();
