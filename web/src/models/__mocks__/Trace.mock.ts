import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import {TRawTrace, TTrace} from '../../types/Trace.types';
import Trace from '../Trace.model';
import SpanMock from './Span.mock';

const TraceMock: IMockFactory<TTrace, TRawTrace> = () => ({
  raw(data = {}) {
    return {
      traceId: faker.datatype.uuid(),
      tree: SpanMock.raw(),
      flat: {
        '1': SpanMock.raw(),
        '2': SpanMock.raw(),
      },
      ...data,
    };
  },
  model(data = {}) {
    return Trace(this.raw(data));
  },
});

export default TraceMock();
