import {createSelector} from '@reduxjs/toolkit';

import {endpoints} from 'redux/apis/TraceTest.api';
import {RootState} from '../redux/store';
import AssertionService from '../services/Assertion.service';
import {TSpanSelector} from '../types/Assertion.types';

const stateSelector = (state: RootState) => state;

const AssertionSelectors = () => ({
  selectAffectedSpanList(testId: string, runId: string, selectorList: TSpanSelector[]) {
    return createSelector(stateSelector, state => {
      const query = AssertionService.getSelectorString(selectorList);
      const {data: selectedSpanIdList} = endpoints.getSelectedSpans.select({testId, runId, query})(state);
      const {data: {trace} = {}} = endpoints.getRunById.select({testId, runId})(state);

      if (!selectedSpanIdList) return [];

      return trace?.spans.filter(({id}) => selectedSpanIdList.includes(id)) || [];
    });
  },
  selectAttributeList(testId: string, runId: string, selectorList: TSpanSelector[]) {
    return createSelector(this.selectAffectedSpanList(testId, runId, selectorList), spanList =>
      spanList.flatMap(span => span.attributeList)
    );
  },
});

export default AssertionSelectors();
