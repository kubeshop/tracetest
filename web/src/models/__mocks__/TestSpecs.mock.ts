import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import {TRawTestSpecs, TTestSpecs} from '../../types/TestSpecs.types';
import TestDefinition from '../TestSpecs.model';
import AssertionMock from './Assertion.mock';

const TestSpecsMock: IMockFactory<TTestSpecs, TRawTestSpecs> = () => ({
  raw(data = {}) {
    return {
      specs: new Array(faker.datatype.number({min: 2, max: 10})).fill(null).map((item, index) => ({
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

export default TestSpecsMock();
