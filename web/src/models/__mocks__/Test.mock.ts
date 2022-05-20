import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import {TRawTest, TTest} from '../../types/Test.types';
import Test from '../Test.model';

const TestMock: IMockFactory<TTest, TRawTest> = () => ({
  raw(data = {}) {
    return {
      id: faker.datatype.uuid(),
      name: faker.name.firstName(),
      definition: undefined,
      ...data,
    };
  },
  model(data = {}) {
    return Test(this.raw(data));
  },
});

export default TestMock();
