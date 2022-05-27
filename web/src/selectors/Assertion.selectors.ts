import {createSelector} from '@reduxjs/toolkit';

import {endpoints} from 'redux/apis/TraceTest.api';
import {RootState} from '../redux/store';

const stateSelector = (state: RootState) => state;
const paramsSelector = (state: RootState, testId: string, runId: string, spanIdList: string[]) => ({
  spanIdList,
  testId,
  runId,
});

const selectAffectedSpanList = createSelector(stateSelector, paramsSelector, (state, {spanIdList, testId, runId}) => {
  if (!spanIdList) return [];
  const {data: {trace} = {}} = endpoints.getRunById.select({testId, runId})(state);

  return trace?.spans.filter(({id}) => spanIdList.includes(id)) || [];
});

const selectResultAssertionsBySpan = createSelector(
  stateSelector,
  (_: RootState, testId: string, runId: string, spanId: string) => ({testId, runId, spanId}),
  (state, {testId, runId, spanId}) => {
    const {data} = endpoints.getRunById.select({testId, runId})(state);

    console.log('### data', data);
    if (!data) return {};

    // result.resultList = array of selectors - 1 item for each selector
    // result.resultList[].resultList = array of checks/assertions - 1 item for each check/assertion
    const results = data.result.resultList
      .flatMap(selector =>
        selector.resultList.map(assertion => ({
          id: selector.id,
          attribute: assertion.assertion.attribute,
          result: assertion.spanResults.find(spanResult => spanResult.spanId === spanId),
        }))
      )
      // Filter if it has result for the spanId
      .filter(assertion => Boolean(assertion?.result))
      // Hash items by attribute
      .reduce((prev: {[key: string]: any}, curr) => {
        const value = prev[curr.attribute] || {failed: [], passed: []};

        if (curr.result?.passed) value.passed.push({id: curr.id});
        else value.failed.push({id: curr.id});

        return {...prev, [curr.attribute]: value};
      }, {});

    return results;
  }
);

const AssertionSelectors = () => ({
  selectAffectedSpanList,
  selectAttributeList: createSelector(selectAffectedSpanList, spanList => spanList.flatMap(span => span.attributeList)),
  selectResultAssertionsBySpan,
});

export default AssertionSelectors();
