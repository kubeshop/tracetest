import faker from '@faker-js/faker';
import {TAssertionSpanResult, TRawAssertionSpanResult} from '../../types/Assertion.types';
import {IMockFactory} from '../../types/Common.types';
import AssertionSpanResult from '../AssertionSpanResult.model';

const AssertionSpanResultMock: IMockFactory<TAssertionSpanResult, TRawAssertionSpanResult> = () => ({
  raw(data = {}) {
    return {
      spanId: faker.datatype.uuid(),
      observedValue: faker.random.word(),
      passed: faker.datatype.boolean(),
      error: faker.random.word(),
      ...data,
    };
  },
  model(data = {}) {
    return AssertionSpanResult(this.raw(data));
  },
});

export default AssertionSpanResultMock();
