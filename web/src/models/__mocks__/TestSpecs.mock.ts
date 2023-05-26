import faker from '@faker-js/faker';
import {IMockFactory} from 'types/Common.types';
import TestSpecs, {TRawTestSpecs} from '../TestSpecs.model';

const TestSpecsMock: IMockFactory<TestSpecs, TRawTestSpecs> = () => ({
  raw(data = {}) {
    return {
      specs: new Array(faker.datatype.number({min: 2, max: 10})).fill(null).map((item, index) => ({
        selector: `span[http.status_code] = "20${index}"]`,
        selectorParsed: {query: `span[http.status_code] = "20${index}"]`},
        assertions: new Array(faker.datatype.number({min: 2, max: 10}))
          .fill(null)
          .map(() => 'attr:tracetest.span.type = "http"'),
      })),
      ...data,
    };
  },
  model(data = {}) {
    return TestSpecs(this.raw(data));
  },
});

export default TestSpecsMock();
