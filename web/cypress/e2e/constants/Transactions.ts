import {TRawTest} from '../../../src/models/Test.model';
import {POKEMON_HTTP_ENDPOINT} from '../constants/Test';

export const transactionTestList: TRawTest[] = [
  {
    name: 'POST test',
    description: 'transaction',
    serviceUnderTest: {
      triggerType: 'http',
      http: {
        url: `${POKEMON_HTTP_ENDPOINT}/pokemon/import`,
        method: 'POST',
        body: '{"id": 6}',
      },
    },
  },
  {
    name: 'GET test',
    description: 'transaction',
    serviceUnderTest: {
      triggerType: 'http',
      http: {
        url: `${POKEMON_HTTP_ENDPOINT}/pokemon`,
        method: 'GET',
      },
    },
  },
];
