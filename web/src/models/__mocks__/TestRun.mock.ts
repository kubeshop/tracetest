import faker from '@faker-js/faker';
import {TestState} from 'constants/TestRun.constants';
import {IMockFactory} from 'types/Common.types';
import {TRawAssertionResults} from '../AssertionResults.model';
import TestRun, {TRawTestRun} from '../TestRun.model';
import TraceMock from './Trace.mock';

const TestRunMock: IMockFactory<TestRun, TRawTestRun> = () => ({
  raw(data = {}) {
    return {
      id: faker.datatype.number(),
      traceId: faker.datatype.uuid(),
      spanId: faker.datatype.uuid(),
      createdAt: faker.date.past().toISOString(),
      testVersion: faker.datatype.number(),
      completedAt: faker.date.past().toISOString(),
      response: {},
      result: {} as TRawAssertionResults,
      trace: TraceMock.raw(),
      state: TestState.FINISHED,
      ...data,
    };
  },
  model(data = {}) {
    return TestRun(this.raw(data));
  },
});

export default TestRunMock();
