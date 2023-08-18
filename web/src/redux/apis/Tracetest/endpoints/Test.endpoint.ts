import {HTTP_METHOD} from 'constants/Common.constants';
import {SortBy, SortDirection, TracetestApiTags} from 'constants/Test.constants';
import Test, {TRawTest, TRawTestResource, TRawTestResourceList} from 'models/Test.model';
import {PaginationResponse} from 'hooks/usePagination';
import {TTestApiEndpointBuilder} from '../Tracetest.api';

const defaultHeaders = {'content-type': 'application/json', 'X-Tracetest-Augmented': 'true'};

export const testEndpoints = (builder: TTestApiEndpointBuilder) => ({
  createTest: builder.mutation<Test, TRawTestResource>({
    query: newTest => ({
      url: '/tests',
      method: HTTP_METHOD.POST,
      body: newTest,
      headers: defaultHeaders,
    }),
    transformResponse: (rawTest: TRawTestResource) => Test(rawTest),
    invalidatesTags: [
      {type: TracetestApiTags.TEST, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  editTest: builder.mutation<Test, {test: TRawTestResource; testId: string}>({
    query: ({test, testId}) => ({
      url: `/tests/${testId}`,
      method: HTTP_METHOD.PUT,
      body: test,
      headers: defaultHeaders,
    }),
    invalidatesTags: test => [
      {type: TracetestApiTags.TEST, id: 'LIST'},
      {type: TracetestApiTags.TEST, id: test?.id},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  getTestList: builder.query<
    PaginationResponse<Test>,
    {take?: number; skip?: number; query?: string; sortBy?: SortBy; sortDirection?: SortDirection}
  >({
    query: ({take = 25, skip = 0, query = '', sortBy = '', sortDirection = ''}) => ({
      url: `/tests?take=${take}&skip=${skip}&query=${query}&sortBy=${sortBy}&sortDirection=${sortDirection}`,
      headers: defaultHeaders,
    }),
    providesTags: () => [{type: TracetestApiTags.TEST, id: 'LIST'}],
    transformResponse: ({items = [], count = 0}: TRawTestResourceList) => {
      return {
        items: items.map(rawTest => Test(rawTest)),
        total: count,
      };
    },
  }),
  getTestById: builder.query<Test, {testId: string}>({
    query: ({testId}) => ({url: `/tests/${testId}`, headers: defaultHeaders}),
    providesTags: result => [{type: TracetestApiTags.TEST, id: result?.id}],
    transformResponse: (rawTest: TRawTestResource) => Test(rawTest),
  }),
  getTestVersionById: builder.query<Test, {testId: string; version: number}>({
    query: ({testId, version}) => ({url: `/tests/${testId}/version/${version}`, headers: defaultHeaders}),
    providesTags: result => [{type: TracetestApiTags.TEST, id: result?.id}],
    transformResponse: (rawTest: TRawTest) => Test.FromRawTest(rawTest),
    keepUnusedDataFor: 10,
  }),
  deleteTestById: builder.mutation<Test, {testId: string}>({
    query: ({testId}) => ({url: `/tests/${testId}`, method: 'DELETE'}),
    invalidatesTags: [
      {type: TracetestApiTags.TEST, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
});
