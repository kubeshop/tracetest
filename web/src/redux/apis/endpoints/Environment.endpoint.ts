import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import Environment from 'models/Environment.model';
import {TEnvironment, TRawEnvironment} from 'types/Environment.types';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {getTotalCountFromHeaders} from 'utils/Common';

const EnvironmentEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getEnvironments: builder.query<PaginationResponse<TEnvironment>, {take?: number; skip?: number; query?: string}>({
    query: ({take = 25, skip = 0, query = ''}) => `/environments?take=${take}&skip=${skip}&query=${query}`,
    providesTags: () => [{type: TracetestApiTags.ENVIRONMENT, id: 'LIST'}],
    transformResponse: (rawEnvironments: TRawEnvironment[], meta) => ({
      items: rawEnvironments.map(rawEnv => Environment(rawEnv)),
      total: getTotalCountFromHeaders(meta),
    }),
  }),
  createEnvironment: builder.mutation<undefined, TEnvironment>({
    query: environment => ({
      url: '/environments',
      method: HTTP_METHOD.POST,
      body: environment,
    }),
    invalidatesTags: [{type: TracetestApiTags.ENVIRONMENT, id: 'LIST'}],
  }),
  updateEnvironment: builder.mutation<undefined, {environment: TRawEnvironment; environmentId: string}>({
    query: ({environment, environmentId}) => ({
      url: `/environments/${environmentId}`,
      method: HTTP_METHOD.PUT,
      body: environment,
    }),
    invalidatesTags: [{type: TracetestApiTags.ENVIRONMENT, id: 'LIST'}],
  }),
  deleteEnvironment: builder.mutation<undefined, {environmentId: string}>({
    query: ({environmentId}) => ({
      url: `/environments/${environmentId}`,
      method: HTTP_METHOD.DELETE,
    }),
    invalidatesTags: [{type: TracetestApiTags.ENVIRONMENT, id: 'LIST'}],
  }),
});

export default EnvironmentEndpoint;
