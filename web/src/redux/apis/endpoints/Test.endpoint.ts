import {TRawTest, TTest, TTestApiEndpointBuilder} from 'types/Test.types';
import {HTTP_METHOD} from 'constants/Common.constants';
import {SortBy, SortDirection, TracetestApiTags} from 'constants/Test.constants';
import Test from 'models/Test.model';
import {PaginationResponse} from 'hooks/usePagination';
import {TRawTestSpecs} from 'types/TestSpecs.types';

function getTotalCountFromHeaders(meta: any) {
  return Number(meta?.response?.headers.get('x-total-count') || 0);
}

const TestEndpoint = (builder: TTestApiEndpointBuilder) => ({
  createTest: builder.mutation<TTest, TRawTest>({
    query: newTest => ({
      url: '/tests',
      method: HTTP_METHOD.POST,
      body: newTest,
    }),
    transformResponse: (rawTest: TRawTest) => Test(rawTest),
    invalidatesTags: [{type: TracetestApiTags.TEST, id: 'LIST'}],
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
    ],
  }),
  getTestList: builder.query<
    PaginationResponse<TTest>,
    {take?: number; skip?: number; query?: string; sortBy?: SortBy; sortDirection?: SortDirection}
  >({
    query: ({take = 25, skip = 0, query = '', sortBy = '', sortDirection = ''}) =>
      `/tests?take=${take}&skip=${skip}&query=${query}&sortBy=${sortBy}&sortDirection=${sortDirection}`,
    providesTags: () => [{type: TracetestApiTags.TEST, id: 'LIST'}],
    transformResponse: (rawTestList: TTest[], meta) => {
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
  }),
  deleteTestById: builder.mutation<TTest, {testId: string}>({
    query: ({testId}) => ({url: `/tests/${testId}`, method: 'DELETE'}),
    invalidatesTags: [{type: TracetestApiTags.TEST, id: 'LIST'}],
  }),

  setTestDefinition: builder.mutation<string[], {testId: string; testDefinition: Partial<TRawTestSpecs>}>({
    query: ({testId, testDefinition}) => ({
      url: `/tests/${testId}/definition`,
      method: HTTP_METHOD.PUT,
      body: testDefinition,
    }),
    invalidatesTags: (result, error, {testId}) => [
      {type: TracetestApiTags.TEST, id: testId},
      {type: TracetestApiTags.TEST_DEFINITION, id: testId},
    ],
  }),
});

export default TestEndpoint;
