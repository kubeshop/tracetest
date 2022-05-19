import {createSelector} from '@reduxjs/toolkit';

import {endpoints} from 'redux/apis/Test.api';
import {RootState} from '../redux/store';
import AssertionService from '../services/Assertion.service';
import {IAssertion, IAssertionResultList, IItemSelector, ISpanAssertionResult} from '../types/Assertion.types';

const stateSelector = (state: RootState) => state;

const AssertionSelectors = () => ({
  selectAssertionResultListBySpan(testId = '', resultId = '', spanId = '') {
    return createSelector(stateSelector, (state): IAssertionResultList[] => {
      const {data: test} = endpoints.getTestById.select(testId)(state);
      const {data: result} = endpoints.getResultById.select({testId, resultId})(state);

      if (!spanId || !result) return [];
      const {trace} = result;

      const span = trace?.spans.find(({spanId: id}) => id === spanId);

      return (
        test?.assertions?.reduce<Array<{assertion: IAssertion; assertionResultList: Array<ISpanAssertionResult>}>>(
          (resultList, assertion) => {
            const assertionResultList = AssertionService.runBySpan(span!, assertion);

            return assertionResultList.length ? [...resultList, {assertion, assertionResultList}] : resultList;
          },
          []
        ) || []
      );
    });
  },
  selectAffectedSpanList(testId: string, resultId: string, selectorList: IItemSelector[]) {
    return createSelector(stateSelector, state => {
      const {data: result} = endpoints.getResultById.select({testId, resultId})(state);

      if (!result) return [];
      const {trace} = result;

      return AssertionService.getEffectedSpansCount(trace!, selectorList);
    });
  },
  selectAttributeList(testId: string, resultId: string, selectorList: IItemSelector[]) {
    return createSelector(this.selectAffectedSpanList(testId, resultId, selectorList), spanList =>
      spanList.flatMap(span => span.attributeList)
    );
  },
});

export default AssertionSelectors();
