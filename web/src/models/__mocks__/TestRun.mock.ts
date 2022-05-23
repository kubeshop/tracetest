import faker from '@faker-js/faker';
import {TestState} from '../../constants/TestRun.constants';
import {TRawAssertionResults} from '../../types/Assertion.types';
import {IMockFactory} from '../../types/Common.types';
import {TRawTestRun, TTestRun} from '../../types/TestRun.types';
import TestRunResult from '../TestRun.model';
import TraceMock from './Trace.mock';

const TestRunMock: IMockFactory<TTestRun, TRawTestRun> = () => ({
  raw(data = {}) {
    return {
      id: faker.datatype.uuid(),
      traceId: faker.datatype.uuid(),
      spanId: faker.datatype.uuid(),
      createdAt: faker.date.past().toISOString(),
      completedAt: faker.date.past().toISOString(),
      response: {},
      result: {} as TRawAssertionResults,
      trace: TraceMock.raw(),
      state: TestState.FINISHED,
      ...data,
    };
  },
  model(data = {}) {
    return TestRunResult(this.raw(data));
  },
});

export default TestRunMock();
