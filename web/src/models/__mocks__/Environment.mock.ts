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
      variables: [],
      ...data,
    };
  },
  model(data = {}) {
    return Environment(this.raw(data));
  },
});

export default EnvironmentMock();
