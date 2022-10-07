import faker from '@faker-js/faker';
import {IEnvironment} from '../../pages/Environments/IEnvironment';
import {IMockFactory} from '../../types/Common.types';
import Environment from '../Environment.model';

const EnvironmentMock: IMockFactory<IEnvironment, IEnvironment> = () => ({
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
