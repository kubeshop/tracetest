import {HTTP_METHOD} from 'constants/Common.constants';
import {SortBy, SortDirection, TracetestApiTags} from 'constants/Test.constants';
import Test from 'models/Test.model';
import {PaginationResponse} from 'hooks/usePagination';
import {TRawTest, TTest, TTestApiEndpointBuilder} from 'types/Test.types';
import {TRawTestOutput} from 'types/TestOutput.types';
import {TRawTestSpecs} from 'types/TestSpecs.types';
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
  setTestOutputs: builder.mutation<undefined, {testId: string; testOutputs: TRawTestOutput[]}>({
    query: ({testId, testOutputs}) => ({
      url: `/tests/${testId}/outputs`,
      method: HTTP_METHOD.PUT,
      body: testOutputs,
    }),
    invalidatesTags: (result, error, {testId}) => [{type: TracetestApiTags.TEST, id: testId}],
  }),
});

export default TestEndpoint;
