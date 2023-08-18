import {HTTP_METHOD} from 'constants/Common.constants';
import {SortBy, SortDirection, TracetestApiTags} from 'constants/Test.constants';
import TestSuite, {TRawTestSuite, TRawTestSuiteResource, TRawTestSuiteResourceList} from 'models/TestSuite.model';
import TestSuiteService from 'services/TestSuite.service';
import {TDraftTestSuite} from 'types/TestSuite.types';
import {PaginationResponse} from 'hooks/usePagination';
import {TTestApiEndpointBuilder} from '../Tracetest.api';

const defaultHeaders = {'content-type': 'application/json', 'X-Tracetest-Augmented': 'true'};

export const testSuiteEndpoints = (builder: TTestApiEndpointBuilder) => ({
  createTestSuite: builder.mutation<TestSuite, TDraftTestSuite>({
    query: suite => ({
      url: '/testsuites',
      method: HTTP_METHOD.POST,
      body: TestSuiteService.getRawFromDraft(suite),
      headers: defaultHeaders,
    }),
    transformResponse: (raw: TRawTestSuiteResource) => TestSuite(raw),
    invalidatesTags: [
      {type: TracetestApiTags.TESTSUITE, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  editTestSuite: builder.mutation<TestSuite, {draft: TDraftTestSuite; testSuiteId: string}>({
    query: ({testSuiteId, draft}) => ({
      url: `/testsuites/${testSuiteId}`,
      method: HTTP_METHOD.PUT,
      body: TestSuiteService.getRawFromDraft(draft),
      headers: defaultHeaders,
    }),
    invalidatesTags: [{type: TracetestApiTags.TESTSUITE, id: 'LIST'}],
  }),
  deleteTestSuiteById: builder.mutation<TestSuite, {testSuiteId: string}>({
    query: ({testSuiteId}) => ({
      url: `/testsuites/${testSuiteId}`,
      method: HTTP_METHOD.DELETE,
      headers: defaultHeaders,
    }),
    invalidatesTags: [
      {type: TracetestApiTags.TESTSUITE, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  getTestSuiteById: builder.query<TestSuite, {testSuiteId: string}>({
    query: ({testSuiteId}) => ({
      url: `/testsuites/${testSuiteId}`,
      headers: defaultHeaders,
    }),
    providesTags: result => [{type: TracetestApiTags.TESTSUITE, id: result?.id}],
    transformResponse: (raw: TRawTestSuiteResource) => TestSuite(raw),
  }),
  getTestSuiteList: builder.query<
    PaginationResponse<TestSuite>,
    {take?: number; skip?: number; query?: string; sortBy?: SortBy; sortDirection?: SortDirection}
  >({
    query: ({take = 25, skip = 0, query = '', sortBy = '', sortDirection = ''}) => ({
      url: `/testsuites?take=${take}&skip=${skip}&query=${query}&sortBy=${sortBy}&sortDirection=${sortDirection}`,
      headers: defaultHeaders,
    }),
    providesTags: () => [{type: TracetestApiTags.TESTSUITE, id: 'LIST'}],
    transformResponse: ({items = [], count = 0}: TRawTestSuiteResourceList) => {
      return {
        items: items.map(raw => TestSuite(raw)),
        total: count,
      };
    },
  }),
  getTestSuiteVersionById: builder.query<TestSuite, {testSuiteId: string; version: number}>({
    query: ({testSuiteId, version}) => `/testsuites/${testSuiteId}/version/${version}`,
    providesTags: result => [{type: TracetestApiTags.TESTSUITE, id: `${result?.id}-${result?.version}`}],
    transformResponse: (rawTest: TRawTestSuite) => TestSuite.FromRawTestSuite(rawTest),
    keepUnusedDataFor: 10,
  }),
});
