import {createSelector} from '@reduxjs/toolkit';
import {sortBy} from 'lodash';

import TracetestAPI from 'redux/apis/Tracetest';
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

const attributeKeySelector = (state: RootState, testId: string, runId: string, spanIdList: string[], key: string) =>
  key;

const selectMatchedSpanList = createSelector(stateSelector, paramsSelector, (state, {spanIdList, testId, runId}) => {
  const {data: {trace} = {}} = TracetestAPI.instance.endpoints.getRunById.select({testId, runId})(state);
  if (!spanIdList.length) return trace?.spans || [];

  return trace?.spans.filter(({id}) => spanIdList.includes(id)) || [];
});

const AssertionSelectors = () => {
  return {
    selectMatchedSpanList,
    selectAttributeList: createSelector(
      selectMatchedSpanList,
      SpanSelectors.selectMatchedSpans,
      (spanList, matchedSpans) =>
        spanList
          .flatMap(span => span.attributeList)
          .concat(SpanAttributeService.getPseudoAttributeList(matchedSpans.length))
    ),
    selectAllAttributeList: createSelector(stateSelector, paramsSelector, (state, {testId, runId}) => {
      const {data: {trace} = {}} = TracetestAPI.instance.endpoints.getRunById.select({testId, runId})(state);

      const spanList = trace?.spans || [];

      return spanList.flatMap(span => span.attributeList);
    }),
    selectSelectorAttributeList: createSelector(
      selectMatchedSpanList,
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
    selectAttributeValueList: createSelector(selectMatchedSpanList, attributeKeySelector, (spanList, attrKey) => {
      const list = spanList.reduce<string[]>((acc, span) => {
        const attr = span.attributeList.find(({key}) => key === attrKey);

        if (attr) return acc.concat(attr.value);
        return acc;
      }, []);

      return list;
    }),
  };
};

export default AssertionSelectors();
