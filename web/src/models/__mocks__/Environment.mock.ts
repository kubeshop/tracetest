import faker from '@faker-js/faker';
import {TEnvironment} from '../../types/Environment.types';
import {IMockFactory} from '../../types/Common.types';
import Environment from '../Environment.model';

const EnvironmentMock: IMockFactory<TEnvironment, TEnvironment> = () => ({
  raw(data = {}) {
    return {
      id: faker.datatype.uuid(),
      name: faker.name.jobType(),
      description: faker.name.jobDescriptor(),
      variables: [
        {
          key: 'HOST',
          value: 'http://localhost',
        },
        {
          key: 'PORT',
          value: '3000',
        },
        {
          key: 'AUTH_TOKEN',
          value: '12313215akashdlasjkldql;dqkwl;wm',
        },
      ],
      ...data,
    };
  },
  model(data = {}) {
    return Environment(this.raw(data));
  },
});

export default EnvironmentMock();
