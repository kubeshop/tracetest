import {createSelector} from '@reduxjs/toolkit';

import {endpoints} from 'redux/apis/TraceTest.api';
import {RootState} from 'redux/store';
import SpanAttributeService from '../services/SpanAttribute.service';
import {TSpanSelector} from '../types/Assertion.types';
import {selectAffectedSpans} from './Span.selectors';

const stateSelector = (state: RootState) => state;
const paramsSelector = (state: RootState, testId: string, runId: string, spanIdList: string[]) => ({
  spanIdList,
  testId,
  runId,
});

const currentSelectorListSelector = (
  state: RootState,
  testId: string,
  runId: string,
  spanIdList: string[],
  currentSelectorList: TSpanSelector[]
) => currentSelectorList.map(({key}) => key);

const selectAffectedSpanList = createSelector(stateSelector, paramsSelector, (state, {spanIdList, testId, runId}) => {
  if (!spanIdList) return [];
  const {data: {trace} = {}} = endpoints.getRunById.select({testId, runId})(state);

  return trace?.spans.filter(({id}) => spanIdList.includes(id)) || [];
});

const AssertionSelectors = () => {
  return {
    selectAffectedSpanList,
    selectAttributeList: createSelector(selectAffectedSpanList, selectAffectedSpans, (spanList, affectedSpans) => spanList.flatMap(span => span.attributeList).concat(SpanAttributeService.additionalAttributes(affectedSpans.length))),
    selectSelectorAttributeList: createSelector(
      selectAffectedSpanList,
      currentSelectorListSelector,
      (spanList, currentSelectorList) =>
        SpanAttributeService.getFilteredSelectorAttributeList(
          spanList.flatMap(span => span.attributeList),
          currentSelectorList
        )
    ),
  };
};

export default AssertionSelectors();
