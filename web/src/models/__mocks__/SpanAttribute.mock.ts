import faker from '@faker-js/faker';
import {IMockFactory} from '../../types/Common.types';
import {IRawSpanAttribute, ISpanAttribute} from '../../types/SpanAttribute.types';
import SpanAttribute from '../SpanAttribute.model';

const SpanAttributeMock: IMockFactory<ISpanAttribute, IRawSpanAttribute> = () => ({
  raw(data = {}) {
    return {
      key: `${faker.random.word()}.${faker.random.word()}`,
      value: {
        stringValue: faker.random.word(),
        intValue: faker.datatype.number(),
        booleanValue: faker.datatype.boolean(),
        doubleValue: faker.datatype.number(),
        kvlistValue: {values: []},
      },
      ...data,
    };
  },
  model(data = {}) {
    return SpanAttribute(this.raw(data));
  },
});

export default SpanAttributeMock();
