import uniq from 'lodash/uniq';

import {TRawAssertionResult} from 'types/Assertion.types';

const AssertionService = () => ({
  getSpanIds(resultList: TRawAssertionResult[]) {
    const spanIds = resultList.flatMap(assertion => assertion?.spanResults?.map(span => span.spanId ?? '') ?? []);
    return uniq(spanIds);
  },
});

export default AssertionService();
