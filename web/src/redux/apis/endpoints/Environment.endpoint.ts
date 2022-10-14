import {HTTP_METHOD} from 'constants/Common.constants';
import {PaginationResponse} from 'hooks/usePagination';
import Environment from 'models/__mocks__/Environment.mock';
import KeyValueMock from 'models/__mocks__/KeyValue.mock';
import {TEnvironment} from 'types/Environment.types';
import {IKeyValue, TracetestApiTags} from 'constants/Test.constants';
import {TTestApiEndpointBuilder} from 'types/Test.types';

const environmentList = [
  Environment.model({id: 'ae7162b3-54e0-4603-9d33-12345', name: 'Production', description: 'Production environment'}),
  Environment.model({
    id: 'ae7162b3-54e0-4603-9d33-423b12cf67c8',
    name: 'Development',
    description: 'Developing environment',
  }),
];

const keyValueListOne = [KeyValueMock.model()];
const keyValueListTwo = [
  KeyValueMock.model({key: 'user', value: 'testAdmin'}),
  KeyValueMock.model({key: 'password', value: '1234'}),
];

const EnvironmentEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getEnvList: builder.query<PaginationResponse<TEnvironment>, {take?: number; skip?: number; query?: string}>({
    query: ({take = 25, skip = 0, query = ''}) => `/tests?take=${take}&skip=${skip}&query=${query}`,
    providesTags: () => [{type: TracetestApiTags.ENVIRONMENT, id: 'LIST'}],
    transformResponse: () => {
      return {
        total: environmentList.length,
        items: environmentList,
      };
    },
  }),
  getEnvironmentSecretList: builder.query<IKeyValue[], {environmentId: string; take?: number; skip?: number}>({
    query: ({take = 25, skip = 0}) => `/tests?take=${take}&skip=${skip}`,
    providesTags: (result, error, {environmentId}) => [
      {type: TracetestApiTags.ENVIRONMENT, id: `${environmentId}-LIST`},
    ],
    transformResponse: (raw, meta, args) => {
      return args.environmentId === 'ae7162b3-54e0-4603-9d33-423b12cf67c8' ? keyValueListOne : keyValueListTwo;
    },
  }),
  createEnvironment: builder.mutation<undefined, TEnvironment>({
    query: newEnvironment => ({
      url: '/environments',
      method: HTTP_METHOD.POST,
      body: newEnvironment,
    }),
    transformResponse: () => undefined,
    invalidatesTags: [{type: TracetestApiTags.ENVIRONMENT, id: 'LIST'}],
  }),
});

export default EnvironmentEndpoint;
