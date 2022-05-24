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
      definition: {
        definitions: [
          {
            selector: faker.random.word(),
            assertionList: [AssertionResultMock.raw()],
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
