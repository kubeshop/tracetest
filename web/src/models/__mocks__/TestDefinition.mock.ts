import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import {TRawTestDefinition, TTestDefinition} from '../../types/TestDefinition.types';
import TestDefinition from '../TestDefinition.model';
import AssertionMock from './Assertion.mock';

const TestDefinitionMock: IMockFactory<TTestDefinition, TRawTestDefinition> = () => ({
  raw(data = {}) {
    return {
      definitions: new Array(faker.datatype.number({min: 2, max: 10})).fill(null).map((item, index) => ({
        selector: {query: `span[http.status_code] = "20${index}"]`},
        assertions: new Array(faker.datatype.number({min: 2, max: 10})).fill(null).map(() => AssertionMock.raw()),
      })),
      ...data,
    };
  },
  model(data = {}) {
    return TestDefinition(this.raw(data));
  },
});

export default TestDefinitionMock();
