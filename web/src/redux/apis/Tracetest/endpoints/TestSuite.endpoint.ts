import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import TestSuite, {TRawTestSuite, TRawTestSuiteResource} from 'models/TestSuite.model';
import TestSuiteService from 'services/TestSuite.service';
import {TDraftTestSuite} from 'types/TestSuite.types';
import TraceTestAPI from '../Tracetest.api';

const defaultHeaders = {'content-type': 'application/json', 'X-Tracetest-Augmented': 'true'};

const testSuiteEndpoints = TraceTestAPI.injectEndpoints({
  endpoints: builder => ({
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
    getTestSuiteVersionById: builder.query<TestSuite, {testSuiteId: string; version: number}>({
      query: ({testSuiteId, version}) => `/testsuites/${testSuiteId}/version/${version}`,
      providesTags: result => [{type: TracetestApiTags.TESTSUITE, id: `${result?.id}-${result?.version}`}],
      transformResponse: (rawTest: TRawTestSuite) => TestSuite.FromRawTestSuite(rawTest),
      keepUnusedDataFor: 10,
    }),
  }),
});

export const {
  useCreateTestSuiteMutation,
  useDeleteTestSuiteByIdMutation,
  useEditTestSuiteMutation,
  useGetTestSuiteByIdQuery,
  useGetTestSuiteVersionByIdQuery,
  useLazyGetTestSuiteByIdQuery,
  useLazyGetTestSuiteVersionByIdQuery,
} = testSuiteEndpoints;
