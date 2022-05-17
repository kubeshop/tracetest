import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import {IRawSpan, ISpan} from '../../types/Span.types';
import Span from '../Span.model';

const SpanMock: IMockFactory<ISpan, IRawSpan> = () => ({
  raw(data = {}) {
    return {
      traceId: faker.datatype.uuid(),
      spanId: faker.datatype.uuid(),
      parentSpanId: faker.datatype.uuid(),
      name: faker.random.word(),
      kind: faker.random.word(),
      startTimeUnixNano: faker.date.recent().toISOString(),
      endTimeUnixNano: faker.date.recent().toISOString(),
      attributes: faker.datatype
        .array(faker.datatype.number({min: 2, max: 10}))
        .map(() => ({
          key: faker.random.word(),
          value: {
            stringValue: faker.random.word(),
            kvlistValue: {values: []},
          },
        }))
        .concat({
          key: 'service.name',
          value: {
            stringValue: 'mock',
            kvlistValue: {values: []},
          },
        }),
      status: {
        code: faker.random.word(),
      },
      ...data,
    };
  },
  model(data = {}) {
    return Span(this.raw(data));
  },
});

export default SpanMock();
