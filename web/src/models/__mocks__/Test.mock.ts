import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import {TRawTest, TTest} from '../../types/Test.types';
import Test from '../Test.model';
import AssertionResultMock from './AssertionResult.mock';

const TestMock: IMockFactory<TTest, TRawTest> = () => ({
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
    return Test(this.raw(data));
  },
});

export default TestMock();
