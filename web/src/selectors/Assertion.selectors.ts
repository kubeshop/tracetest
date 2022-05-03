import {createSelector} from '@reduxjs/toolkit';

import {endpoints} from 'redux/apis/Test.api';
import {RootState} from '../redux/store';
import AssertionService from '../services/Assertion.service';
import {TAssertion, TItemSelector, TSpanAssertionResult} from '../types/Assertion.types';

const stateSelector = (state: RootState) => state;

const AssertionSelectors = () => ({
  selectAssertionResultListBySpan(testId = '', resultId = '', spanId = '') {
    return createSelector(stateSelector, state => {
      const {data: test} = endpoints.getTestById.select(testId)(state);
      const {data: result} = endpoints.getResultById.select({testId, resultId})(state);

      if (!spanId || !result) return [];
      const {trace} = result;

      const span = trace?.spans.find(({spanId: id}) => id === spanId);

      return (
        test?.assertions?.reduce<Array<{assertion: TAssertion; assertionResultList: Array<TSpanAssertionResult>}>>(
          (resultList, assertion) => {
            const assertionResultList = AssertionService.runBySpan(span!, assertion);

            return assertionResultList.length ? [...resultList, {assertion, assertionResultList}] : resultList;
          },
          []
        ) || []
      );
    });
  },
  selectAffectedSpanCount(testId: string, resultId: string, selectorList: TItemSelector[]) {
    return createSelector(stateSelector, state => {
      const {data: result} = endpoints.getResultById.select({testId, resultId})(state);

      if (!result) return 0;
      const {trace} = result;

      return AssertionService.getEffectedSpansCount(trace!, selectorList);
    });
  },
});

export default AssertionSelectors();
