import countBy from 'lodash/countBy';
import uniq from 'lodash/uniq';

import {TAssertionResult, ICheckResult, TRawAssertionResult} from 'types/Assertion.types';

const AssertionService = () => ({
  getSpanIds(resultList: TRawAssertionResult[]) {
    const spanIds = resultList.flatMap(assertion => assertion?.spanResults?.map(span => span.spanId ?? '') ?? []).filter(spanId => Boolean(spanId));
    return uniq(spanIds);
  },

  getTotalPassedChecks(resultList: TAssertionResult[]) {
    const passedResults = resultList.flatMap(({spanResults}) => spanResults.map(({passed}) => passed));
    return countBy(passedResults);
  },

  getResultsHashedBySpanId(resultList: TAssertionResult[]) {
    return resultList
      .flatMap(({assertion, spanResults}) => spanResults.map(spanResult => ({result: spanResult, assertion})))
      .reduce((prev: Record<string, ICheckResult[]>, curr) => {
        const items = prev[curr.result.spanId] || [];
        items.push(curr);

        return {
          ...prev,
          [curr.result.spanId]: items,
        };
      }, {});
  },
});

export default AssertionService();
