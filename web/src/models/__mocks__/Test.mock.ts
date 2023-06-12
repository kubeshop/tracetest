import faker from '@faker-js/faker';
import {IMockFactory} from 'types/Common.types';
import Test, {TRawTest} from '../Test.model';
import AssertionResultMock from './AssertionResult.mock';

const TestMock: IMockFactory<Test, TRawTest> = () => ({
  raw(data = {}) {
    return {
      id: faker.datatype.uuid(),
      name: faker.name.firstName(),
      version: faker.datatype.number(),
      definition: {
        definitions: [
          {
            selector: {query: faker.random.word()},
            assertions: [AssertionResultMock.raw()],
          },
        ],
      },
      ...data,
    };
  },
  model(data = {}) {
    return Test.FromRawTest(this.raw(data));
  },
});

export default TestMock();
