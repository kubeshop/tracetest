import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import Environment, {TRawEnvironment} from 'models/Environment.model';
import {TListResponse} from 'types/Settings.types';
import TraceTestAPI from '../Tracetest.api';

const environmentEndpoints = TraceTestAPI.injectEndpoints({
  endpoints: builder => ({
    getEnvironments: builder.query<PaginationResponse<Environment>, {take?: number; skip?: number; query?: string}>({
      query: ({take = 25, skip = 0, query = ''}) => ({
        url: `/environments?take=${take}&skip=${skip}&query=${query}`,
        headers: {
          'content-type': 'application/json',
        },
      }),
      providesTags: () => [{type: TracetestApiTags.ENVIRONMENT, id: 'LIST'}],
      transformResponse: ({items, count}: TListResponse<TRawEnvironment>) => ({
        items: items.map(rawEnv => Environment(rawEnv)),
        total: count,
      }),
    }),
    createEnvironment: builder.mutation<undefined, TRawEnvironment>({
      query: environment => ({
        url: '/environments',
        method: HTTP_METHOD.POST,
        body: environment,
        headers: {
          'content-type': 'application/json',
        },
      }),
      invalidatesTags: [{type: TracetestApiTags.ENVIRONMENT, id: 'LIST'}],
    }),
    updateEnvironment: builder.mutation<undefined, {environment: TRawEnvironment; environmentId: string}>({
      query: ({environment, environmentId}) => ({
        url: `/environments/${environmentId}`,
        method: HTTP_METHOD.PUT,
        body: environment,
        headers: {
          'content-type': 'application/json',
        },
      }),
      invalidatesTags: [{type: TracetestApiTags.ENVIRONMENT, id: 'LIST'}],
    }),
    deleteEnvironment: builder.mutation<undefined, {environmentId: string}>({
      query: ({environmentId}) => ({
        url: `/environments/${environmentId}`,
        method: HTTP_METHOD.DELETE,
        headers: {
          'content-type': 'application/json',
        },
      }),
      invalidatesTags: [{type: TracetestApiTags.ENVIRONMENT, id: 'LIST'}],
    }),
  }),
});

export const {
  useGetEnvironmentsQuery,
  useCreateEnvironmentMutation,
  useUpdateEnvironmentMutation,
  useDeleteEnvironmentMutation,
} = environmentEndpoints;

export const {endpoints} = environmentEndpoints;
