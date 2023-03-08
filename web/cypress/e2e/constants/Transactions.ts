import {TRawTest} from '../../../src/models/Test.model';

export const transactionTestList: TRawTest[] = [
  {
    name: 'POST test',
    description: 'transaction',
    serviceUnderTest: {
      triggerType: 'http',
      triggerSettings: {
        http: {
          url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon/import',
          method: 'POST',
          body: '{"id": 6}',
        },
      },
    },
  },
  {
    name: 'GET test',
    description: 'transaction',
    serviceUnderTest: {
      triggerType: 'http',
      triggerSettings: {
        http: {
          url: 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon',
          method: 'GET',
        },
      },
    },
  },
];
