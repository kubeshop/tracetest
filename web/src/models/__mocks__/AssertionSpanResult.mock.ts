import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import AssertionSpanResult, {TRawAssertionSpanResult} from '../AssertionSpanResult.model';

const AssertionSpanResultMock: IMockFactory<AssertionSpanResult, TRawAssertionSpanResult> = () => ({
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
