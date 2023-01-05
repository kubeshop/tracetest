import {HTTP_METHOD} from 'constants/Common.constants';
import {SortBy, SortDirection, TracetestApiTags} from 'constants/Test.constants';
import Test from 'models/Test.model';
import {PaginationResponse} from 'hooks/usePagination';
import {TRawTest, TTest, TTestApiEndpointBuilder} from 'types/Test.types';
import {TRawTestVariables, TTestVariables} from 'types/Variables.types';
import TestVariables from 'models/TestVariables.model';
import {getTotalCountFromHeaders} from 'utils/Common';

const TestEndpoint = (builder: TTestApiEndpointBuilder) => ({
  createTest: builder.mutation<TTest, TRawTest>({
    query: newTest => ({
      url: '/tests',
      method: HTTP_METHOD.POST,
      body: newTest,
    }),
    transformResponse: (rawTest: TRawTest) => Test(rawTest),
    invalidatesTags: [
      {type: TracetestApiTags.TEST, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  editTest: builder.mutation<TTest, {test: TRawTest; testId: string}>({
    query: ({test, testId}) => ({
      url: `/tests/${testId}`,
      method: HTTP_METHOD.PUT,
      body: test,
    }),
    invalidatesTags: test => [
      {type: TracetestApiTags.TEST, id: 'LIST'},
      {type: TracetestApiTags.TEST, id: test?.id},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  getTestList: builder.query<
    PaginationResponse<TTest>,
    {take?: number; skip?: number; query?: string; sortBy?: SortBy; sortDirection?: SortDirection}
  >({
    query: ({take = 25, skip = 0, query = '', sortBy = '', sortDirection = ''}) =>
      `/tests?take=${take}&skip=${skip}&query=${query}&sortBy=${sortBy}&sortDirection=${sortDirection}`,
    providesTags: () => [{type: TracetestApiTags.TEST, id: 'LIST'}],
    transformResponse: (rawTestList: TRawTest[], meta) => {
      return {
        items: rawTestList.map(rawTest => Test(rawTest)),
        total: getTotalCountFromHeaders(meta),
      };
    },
  }),
  getTestById: builder.query<TTest, {testId: string}>({
    query: ({testId}) => `/tests/${testId}`,
    providesTags: result => [{type: TracetestApiTags.TEST, id: result?.id}],
    transformResponse: (rawTest: TRawTest) => Test(rawTest),
  }),
  getTestVersionById: builder.query<TTest, {testId: string; version: number}>({
    query: ({testId, version}) => `/tests/${testId}/version/${version}`,
    providesTags: result => [{type: TracetestApiTags.TEST, id: result?.id}],
    transformResponse: (rawTest: TRawTest) => Test(rawTest),
    keepUnusedDataFor: 10,
  }),
  deleteTestById: builder.mutation<TTest, {testId: string}>({
    query: ({testId}) => ({url: `/tests/${testId}`, method: 'DELETE'}),
    invalidatesTags: [
      {type: TracetestApiTags.TEST, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  getTestVariables: builder.query<
    TTestVariables,
    {testId: string; version: number; environmentId?: string; runId?: string}
  >({
    query: ({testId, environmentId = '', version, runId = '0'}) =>
      `/tests/${testId}/version/${version}/variables?environmentId=${environmentId}&runId=${runId}`,
    providesTags: (result, error, {testId}) => [{type: TracetestApiTags.TEST, id: `${testId}-variables`}],
    transformResponse: (rawTestVariables: TRawTestVariables) => TestVariables(rawTestVariables),
  }),
});

export default TestEndpoint;
