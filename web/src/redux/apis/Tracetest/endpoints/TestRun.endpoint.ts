import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import AssertionResults, {TRawAssertionResults} from 'models/AssertionResults.model';
import {TVariableSetValue} from 'models/VariableSet.model';
import RunError from 'models/RunError.model';
import SelectedSpans, {TRawSelectedSpans} from 'models/SelectedSpans.model';
import Test from 'models/Test.model';
import TestRun, {TRawTestRun} from 'models/TestRun.model';
import TestRunEvent, {TRawTestRunEvent} from 'models/TestRunEvent.model';
import SearchSpansResult, {TRawSearchSpansResult} from 'models/SearchSpansResult.model';
import {KnownSources} from 'models/RunMetadata.model';
import {TRawTestSpecs} from 'models/TestSpecs.model';
import {TTestApiEndpointBuilder} from '../Tracetest.api';

function getTotalCountFromHeaders(meta: any) {
  return Number(meta?.response?.headers.get('x-total-count') || 0);
}

export const testRunEndpoints = (builder: TTestApiEndpointBuilder) => ({
  runTest: builder.mutation<TestRun, {testId: string; variableSetId?: string; variables?: TVariableSetValue[]}>({
    query: ({testId, variableSetId, variables = []}) => ({
      url: `/tests/${testId}/run`,
      method: HTTP_METHOD.POST,
      body: {
        variableSetId,
        variables,
        metadata: {
          source: KnownSources.WEB,
        },
      },
    }),
    invalidatesTags: (response, error, {testId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-LIST`},
      {type: TracetestApiTags.TEST, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
      {type: TracetestApiTags.TEST, id: testId},
    ],
    transformResponse: (rawTestRun: TRawTestRun) => TestRun(rawTestRun),
    transformErrorResponse: ({data: result}) => RunError(result),
  }),
  getRunList: builder.query<PaginationResponse<TestRun>, {testId: string; take?: number; skip?: number}>({
    query: ({testId, take = 25, skip = 0}) => `/tests/${testId}/run?take=${take}&skip=${skip}`,
    providesTags: (result, error, {testId}) => [{type: TracetestApiTags.TEST_RUN, id: `${testId}-LIST`}],
    transformResponse: (rawTestResultList: TRawTestRun[], meta) => ({
      total: getTotalCountFromHeaders(meta),
      items: rawTestResultList.map(rawTestResult => TestRun(rawTestResult)),
    }),
  }),
  getRunById: builder.query<TestRun, {runId: number; testId: string}>({
    query: ({testId, runId}) => `/tests/${testId}/run/${runId}`,
    providesTags: result => (result ? [{type: TracetestApiTags.TEST_RUN, id: result?.id}] : []),
    transformResponse: (rawTestResult: TRawTestRun) => TestRun(rawTestResult),
  }),
  reRun: builder.mutation<TestRun, {testId: string; runId: number}>({
    query: ({testId, runId}) => ({
      url: `/tests/${testId}/run/${runId}/rerun`,
      method: HTTP_METHOD.POST,
    }),
    invalidatesTags: (response, error, {testId, runId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-LIST`},
      {type: TracetestApiTags.TEST_RUN, id: runId},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
      {type: TracetestApiTags.TEST, id: testId},
    ],
    transformResponse: (rawTestRun: TRawTestRun) => TestRun(rawTestRun),
  }),
  dryRun: builder.mutation<AssertionResults, {testId: string; runId: number; testDefinition: Partial<TRawTestSpecs>}>({
    query: ({testId, runId, testDefinition}) => ({
      url: `/tests/${testId}/run/${runId}/dry-run`,
      method: HTTP_METHOD.PUT,
      body: testDefinition,
    }),
    transformResponse: (rawTestResults: TRawAssertionResults) => AssertionResults(rawTestResults),
  }),
  deleteRunById: builder.mutation<Test, {testId: string; runId: number}>({
    query: ({testId, runId}) => ({url: `/tests/${testId}/run/${runId}`, method: 'DELETE'}),
    invalidatesTags: (result, error, {testId}) => [
      {type: TracetestApiTags.TEST_RUN},
      {type: TracetestApiTags.TEST, id: `${testId}-LIST`},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  stopRun: builder.mutation<null, {runId: number; testId: string}>({
    query: ({runId, testId}) => ({
      url: `/tests/${testId}/run/${runId}/stop`,
      method: HTTP_METHOD.POST,
    }),
    invalidatesTags: (response, error, {testId, runId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-LIST`},
      {type: TracetestApiTags.TEST_RUN, id: runId},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  skipPolling: builder.mutation<null, {runId: number; testId: string}>({
    query: ({runId, testId}) => ({
      url: `/tests/${testId}/run/${runId}/skipPolling`,
      method: HTTP_METHOD.POST,
    }),
    invalidatesTags: (response, error, {testId, runId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-LIST`},
      {type: TracetestApiTags.TEST_RUN, id: runId},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  getJUnitByRunId: builder.query<string, {testId: string; runId: number}>({
    query: ({testId, runId}) => ({url: `/tests/${testId}/run/${runId}/junit.xml`, responseHandler: 'text'}),
    providesTags: (result, error, {testId, runId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-${runId}-junit`},
    ],
  }),
  getSelectedSpans: builder.query<SelectedSpans, {testId: string; runId: number; query: string}>({
    query: ({testId, runId, query}) => `/tests/${testId}/run/${runId}/select?query=${encodeURIComponent(query)}`,
    providesTags: (result, error, {query}) => (result ? [{type: TracetestApiTags.SPAN, id: `${query}-LIST`}] : []),
    transformResponse: (rawSpanList: TRawSelectedSpans) => SelectedSpans(rawSpanList),
  }),
  getSearchedSpans: builder.mutation<SearchSpansResult, {query: string; testId: string; runId: number}>({
    query: ({query, testId, runId}) => ({
      url: `/tests/${testId}/run/${runId}/search-spans`,
      method: HTTP_METHOD.POST,
      body: JSON.stringify({query}),
    }),
    transformResponse: (raw: TRawSearchSpansResult) => SearchSpansResult(raw),
  }),

  getRunEvents: builder.query<TestRunEvent[], {runId: number; testId: string}>({
    query: ({runId, testId}) => `/tests/${testId}/run/${runId}/events`,
    providesTags: [{type: TracetestApiTags.TEST_RUN, id: 'EVENTS'}],
    transformResponse: (rawTestRunEvent: TRawTestRunEvent[]) => rawTestRunEvent.map(event => TestRunEvent(event)),
  }),
});
