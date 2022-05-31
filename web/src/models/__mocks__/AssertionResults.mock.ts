import faker from '@faker-js/faker';
import {TAssertionResults, TRawAssertionResults} from '../../types/Assertion.types';
import {IMockFactory} from '../../types/Common.types';
import AssertionResults from '../AssertionResults.model';
import AssertionResultMock from './AssertionResult.mock';

const AssertionResultsMock: IMockFactory<TAssertionResults, TRawAssertionResults> = () => ({
  raw(data = {}) {
    return {
      allPassed: faker.datatype.boolean(),
      results: new Array(faker.datatype.number({min: 2, max: 10})).fill(null).map((item, index) => ({
        selector: `span[http.status_code] = "20${index}"]`,
        results: new Array(faker.datatype.number({min: 2, max: 10})).fill(null).map(() => AssertionResultMock.raw()),
      })),
      ...data,
    };
  },
  model(data = {}) {
    return AssertionResults(this.raw(data));
  },
});

export default AssertionResultsMock();
