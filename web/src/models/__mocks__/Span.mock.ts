import faker from '@faker-js/faker';
import {IMockFactory} from 'types/Common.types';
import Span, {TRawSpan} from '../Span.model';

const SpanMock: IMockFactory<Span, TRawSpan> = () => ({
  raw(data = {}) {
    return {
      id: faker.datatype.uuid(),
      parentId: faker.datatype.uuid(),
      name: faker.random.word(),
      startTime: faker.date.recent().getMilliseconds(),
      endTime: faker.date.recent().getMilliseconds(),
      attributes: {
        'service.name': 'mock',
        'tracetest.span.duration': '10',
      },
      children: [],
      ...data,
    };
  },
  model(data = {}) {
    return Span(this.raw(data));
  },
});

export default SpanMock();
