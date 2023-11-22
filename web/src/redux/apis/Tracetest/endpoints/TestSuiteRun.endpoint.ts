import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import {TVariableSetValue} from 'models/VariableSet.model';
import RunError from 'models/RunError.model';
import TestSuiteRun, {TRawTestSuiteRunResourceRun} from 'models/TestSuiteRun.model';
import {KnownSources} from 'models/RunMetadata.model';
import {getTotalCountFromHeaders} from 'utils/Common';
import {TTestApiEndpointBuilder} from '../Tracetest.api';

export const testSuiteRunEndpoints = (builder: TTestApiEndpointBuilder) => ({
  runTestSuite: builder.mutation<
    TestSuiteRun,
    {testSuiteId: string; variableSetId?: string; variables?: TVariableSetValue[]}
  >({
    query: ({testSuiteId, variableSetId, variables = []}) => ({
      url: `/testsuites/${testSuiteId}/run`,
      method: HTTP_METHOD.POST,
      body: {variableSetId, variables, metadata: {source: KnownSources.WEB}},
    }),
    invalidatesTags: (result, error, {testSuiteId}) => [
      {type: TracetestApiTags.TESTSUITE_RUN, id: `${testSuiteId}-LIST`},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
    transformResponse: (raw: TRawTestSuiteRunResourceRun) => TestSuiteRun(raw),
    transformErrorResponse: ({data: result}) => RunError(result),
  }),

  getTestSuiteRuns: builder.query<
    PaginationResponse<TestSuiteRun>,
    {testSuiteId: string; take?: number; skip?: number}
  >({
    query: ({testSuiteId, take = 25, skip = 0}) => `/testsuites/${testSuiteId}/run?take=${take}&skip=${skip}`,
    providesTags: (result, error, {testSuiteId}) => [
      {type: TracetestApiTags.TESTSUITE_RUN, id: `${testSuiteId}-LIST`},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
    transformResponse: (raw: TRawTestSuiteRunResourceRun[], meta) => ({
      total: getTotalCountFromHeaders(meta),
      items: raw.map(rawRun => TestSuiteRun(rawRun)),
    }),
  }),

  getTestSuiteRunById: builder.query<TestSuiteRun, {testSuiteId: string; runId: number}>({
    query: ({testSuiteId, runId}) => `/testsuites/${testSuiteId}/run/${runId}`,
    providesTags: result => [{type: TracetestApiTags.TESTSUITE_RUN, id: result?.id}],
    transformResponse: (raw: TRawTestSuiteRunResourceRun) => TestSuiteRun(raw),
  }),

  deleteTestSuiteRunById: builder.mutation<TestSuiteRun, {testSuiteId: string; runId: number}>({
    query: ({testSuiteId, runId}) => ({
      url: `/testsuites/${testSuiteId}/run/${runId}`,
      method: HTTP_METHOD.DELETE,
    }),
    invalidatesTags: (result, error, {testSuiteId}) => [
      {type: TracetestApiTags.TESTSUITE_RUN, id: `${testSuiteId}-LIST`},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
});
