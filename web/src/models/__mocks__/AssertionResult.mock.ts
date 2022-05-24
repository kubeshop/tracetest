import faker from '@faker-js/faker';
import {TAssertionResult, TRawAssertionResult} from '../../types/Assertion.types';
import {IMockFactory} from '../../types/Common.types';
import AssertionResult from '../AssertionResult.model';
import AssertionMock from './Assertion.mock';
import AssertionSpanResultMock from './AssertionSpanResult.mock';

const AssertionResultMock: IMockFactory<TAssertionResult, TRawAssertionResult> = () => ({
  raw(data = {}) {
    return {
      allPassed: faker.datatype.boolean(),
      assertion: AssertionMock.raw(),
      spanResults: new Array(4).fill(null).map(() => AssertionSpanResultMock.raw()),
      ...data,
    };
  },
  model(data = {}) {
    return AssertionResult(this.raw(data));
  },
});

export default AssertionResultMock();
