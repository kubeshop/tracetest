import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import {IRawTrace, ITrace} from '../../types/Trace.types';
import Trace from '../Trace.model';
import SpanMock from './Span.mock';

const TraceMock: IMockFactory<ITrace, IRawTrace> = () => ({
  raw(data = {}) {
    return {
      description: faker.random.words(),
      resourceSpans: [
        {
          resource: {
            attributes: [],
          },
          instrumentationLibrarySpans: [
            {
              instrumentationLibrary: {
                version: String(faker.datatype.number()),
                name: faker.random.word(),
              },
              spans: faker.datatype.array(faker.datatype.number({min: 2, max: 10})).map(() => SpanMock.raw()),
            },
          ],
        },
      ],
      ...data,
    };
  },
  model(data = {}) {
    return Trace(this.raw(data));
  },
});

export default TraceMock();
