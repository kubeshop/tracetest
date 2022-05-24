import faker from '@faker-js/faker';
import {TAssertion, TRawAssertion} from '../../types/Assertion.types';
import {IMockFactory} from '../../types/Common.types';
import Assertion from '../Assertion.model';

const operatorSymbolList = ['=', '<', '>', '!=', '>=', '<=', 'contains'];

const AssertionMock: IMockFactory<TAssertion, TRawAssertion> = () => ({
  raw(data = {}) {
    return {
      attribute: faker.random.word(),
      comparator: operatorSymbolList[faker.datatype.number({min: 0, max: operatorSymbolList.length - 1})],
      expected: faker.random.word(),
      ...data,
    };
  },
  model(data = {}) {
    return Assertion(this.raw(data));
  },
});

export default AssertionMock();
