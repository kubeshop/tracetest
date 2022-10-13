import faker from '@faker-js/faker';
import Transaction from 'models/Transaction.model';
import {IMockFactory} from 'types/Common.types';
import {TTransaction} from 'types/Transaction.types';
import TestMock from './Test.mock';

const TransactionMock: IMockFactory<TTransaction, TTransaction> = () => ({
  raw(data = {}) {
    const test = TestMock.model();
    const test2 = TestMock.model();
    const test3 = TestMock.model();
    return {
      id: faker.datatype.uuid(),
      name: faker.name.firstName(),
      version: faker.datatype.number(),
      description: faker.company.catchPhraseDescriptor(),
      steps: [
        {...test, result: 'success'},
        {...test2, result: 'fail'},
        {...test3, result: 'running'},
      ],
      env: {
        usename: 'john doe',
      },
      ...data,
    };
  },
  model(data = {}) {
    return Transaction(this.raw(data));
  },
});

export default TransactionMock();
