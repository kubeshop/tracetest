import faker from '@faker-js/faker';

import TestRunEvent, {TRawTestRunEvent} from 'models/TestRunEvent.model';
import {IMockFactory} from 'types/Common.types';

const TestRunEventMock: IMockFactory<TestRunEvent, TRawTestRunEvent> = () => ({
  raw(data = {}) {
    return {
      type: `${faker.lorem.slug()}_${faker.helpers.arrayElement(['ERROR', 'INFO', 'START', 'SUCCESS', 'WARNING'])}`,
      stage: faker.helpers.arrayElement(['trigger', 'trace', 'test']),
      title: faker.lorem.sentence(),
      description: faker.lorem.lines(),
      createdAt: faker.date.past().toISOString(),
      testId: faker.datatype.uuid(),
      runId: faker.datatype.uuid(),
      ...data,
    };
  },
  model(data = {}) {
    return TestRunEvent(this.raw(data));
  },
});

export default TestRunEventMock();
