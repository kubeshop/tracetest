import {createSelector} from '@reduxjs/toolkit';

import {endpoints} from 'redux/apis/TraceTest.api';
import {RootState} from 'redux/store';
import SpanAttributeService from '../services/SpanAttribute.service';

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

const AssertionSelectors = () => ({
  selectAffectedSpanList,
  selectAttributeList: createSelector(selectAffectedSpanList, spanList => spanList.flatMap(span => span.attributeList)),
  selectSelectorAttributeList: createSelector(selectAffectedSpanList, spanList =>
    SpanAttributeService.getFilteredSelectorAttributeList(spanList.flatMap(span => span.attributeList))
  ),
});

export default AssertionSelectors();
