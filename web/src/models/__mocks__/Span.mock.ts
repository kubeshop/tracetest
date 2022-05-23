import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import {TRawSpan, TSpan} from '../../types/Span.types';
import Span from '../Span.model';

const SpanMock: IMockFactory<TSpan, TRawSpan> = () => ({
  raw(data = {}) {
    return {
      id: faker.datatype.uuid(),
      parentId: faker.datatype.uuid(),
      name: faker.random.word(),
      startTime: faker.date.recent().toISOString(),
      endTime: faker.date.recent().toISOString(),
      attributes: {
        'service.name': 'mock',
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
