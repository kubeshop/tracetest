import {createSelector} from '@reduxjs/toolkit';
import {sortBy} from 'lodash';

import {endpoints} from 'redux/apis/TraceTest.api';
import {RootState} from 'redux/store';
import SpanAttributeService from '../services/SpanAttribute.service';
import {TSpanSelector} from '../types/Assertion.types';
import SpanSelectors from './Span.selectors';

const stateSelector = (state: RootState) => state;
const paramsSelector = (state: RootState, testId: string, runId: string, spanIdList: string[] = []) => ({
  spanIdList,
  testId,
  runId,
});

const currentSelectorListSelector = (
  state: RootState,
  testId: string,
  runId: string,
  spanIdList?: string[],
  currentSelectorList: TSpanSelector[] = []
) => currentSelectorList.map(({key}) => key);

const selectAffectedSpanList = createSelector(stateSelector, paramsSelector, (state, {spanIdList, testId, runId}) => {
  const {data: {trace} = {}} = endpoints.getRunById.select({testId, runId})(state);
  if (!spanIdList.length) return trace?.spans || [];

  return trace?.spans.filter(({id}) => spanIdList.includes(id)) || [];
});

const AssertionSelectors = () => {
  return {
    selectAffectedSpanList,
    selectAttributeList: createSelector(
      selectAffectedSpanList,
      SpanSelectors.selectAffectedSpans,
      (spanList, affectedSpans) =>
        spanList
          .flatMap(span => span.attributeList)
          .concat(SpanAttributeService.getPseudoAttributeList(affectedSpans.length))
    ),
    selectAllAttributeList: createSelector(stateSelector, paramsSelector, (state, {testId, runId}) => {
      const {data: {trace} = {}} = endpoints.getRunById.select({testId, runId})(state);

      const spanList = trace?.spans || [];

      return spanList.flatMap(span => span.attributeList);
    }),
    selectSelectorAttributeList: createSelector(
      selectAffectedSpanList,
      currentSelectorListSelector,
      (spanList, currentSelectorList) =>
        sortBy(
          SpanAttributeService.getFilteredSelectorAttributeList(
            spanList.flatMap(span => span.attributeList),
            currentSelectorList
          ),
          'key'
        )
    ),
  };
};

export default AssertionSelectors();
