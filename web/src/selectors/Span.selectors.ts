import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import {endpoints} from '../redux/apis/Tracetest';

const spansStateSelector = (state: RootState) => state.spans;
const stateSelector = (state: RootState) => state;
const paramsSelector = (state: RootState, spanId: string, testId: string, runId: string) => ({spanId, testId, runId});

const selectMatchedSpans = createSelector(
  spansStateSelector,
  stateSelector,
  ({matchedSpans}, {testSpecs: {assertionResults, selectedSpec}}) => {
    if (!selectedSpec) return matchedSpans;

    const foundAssertion = assertionResults?.resultList.find(({selector}) => selector === selectedSpec);

    return !foundAssertion ? [] : matchedSpans;
  }
);

const SpanSelectors = () => ({
  selectMatchedSpans,
  selectSpanById: createSelector(stateSelector, paramsSelector, (state, {spanId, testId, runId}) => {
    const {data: {trace} = {}} = endpoints.getRunById.select({testId, runId})(state);

    const spanList = trace?.spans || [];

    return spanList.find(span => span.id === spanId);
  }),
  selectSelectedSpan: createSelector(spansStateSelector, ({selectedSpan}) => selectedSpan),
  selectFocusedSpan: createSelector(spansStateSelector, ({focusedSpan}) => focusedSpan),
});

export default SpanSelectors();
