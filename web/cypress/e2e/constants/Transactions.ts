import {TRawTestResource} from '../../../src/models/Test.model';
import {POKEMON_HTTP_ENDPOINT} from '../constants/Test';

export const transactionTestList: TRawTestResource[] = [
  {
    type: 'Test',
    spec: {
      name: 'POST test',
      description: 'transaction',
      trigger: {
        type: 'http',
        triggerType: 'http',
        http: {
          url: `${POKEMON_HTTP_ENDPOINT}/pokemon`,
          method: 'GET',
        },
      },
      serviceUnderTest: {
        type: 'http',
        triggerType: 'http',
        http: {
          url: `${POKEMON_HTTP_ENDPOINT}/pokemon/import`,
          method: 'POST',
          body: '{"id": 6}',
        },
      },
    },
  },
  {
    type: 'Test',
    spec: {
      name: 'GET test',
      description: 'transaction',
      trigger: {
        triggerType: 'http',
        type: 'http',
        http: {
          url: `${POKEMON_HTTP_ENDPOINT}/pokemon`,
          method: 'GET',
        },
      },
      serviceUnderTest: {
        triggerType: 'http',
        type: 'http',
        http: {
          url: `${POKEMON_HTTP_ENDPOINT}/pokemon`,
          method: 'GET',
        },
      },
    },
  },
];
