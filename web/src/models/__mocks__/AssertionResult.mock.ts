import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import AssertionResult, {TRawAssertionResult} from '../AssertionResult.model';
import AssertionSpanResultMock from './AssertionSpanResult.mock';

const AssertionResultMock: IMockFactory<AssertionResult, TRawAssertionResult> = () => ({
  raw(data = {}) {
    return {
      allPassed: faker.datatype.boolean(),
      assertion: 'attr:tracetest.span.type = "http',
      spanResults: new Array(4).fill(null).map(() => AssertionSpanResultMock.raw()),
      ...data,
    };
  },
  model(data = {}) {
    return AssertionResult(this.raw(data));
  },
});

export default AssertionResultMock();
