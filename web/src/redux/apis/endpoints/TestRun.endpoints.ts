import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import AssertionResults, {TRawAssertionResults} from 'models/AssertionResults.model';
import SelectedSpans, {TRawSelectedSpans} from 'models/SelectedSpans.model';
import TestRun, {TRawTestRun} from 'models/TestRun.model';
import WebSocketService, {IListenerFunction} from 'services/WebSocket.service';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import RunError from 'models/RunError.model';
import {TEnvironmentValue} from 'models/Environment.model';
import {TRawTestSpecs} from 'models/TestSpecs.model';
import Test from 'models/Test.model';

function getTotalCountFromHeaders(meta: any) {
  return Number(meta?.response?.headers.get('x-total-count') || 0);
}

const TestRunEndpoint = (builder: TTestApiEndpointBuilder) => ({
  runTest: builder.mutation<TestRun, {testId: string; environmentId?: string; variables?: TEnvironmentValue[]}>({
    query: ({testId, environmentId, variables = []}) => ({
      url: `/tests/${testId}/run`,
      method: HTTP_METHOD.POST,
      body: {
        environmentId,
        variables,
      },
    }),
    invalidatesTags: (response, error, {testId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-LIST`},
      {type: TracetestApiTags.TEST, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
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
  getRunById: builder.query<TestRun, {runId: string; testId: string}>({
    query: ({testId, runId}) => `/tests/${testId}/run/${runId}`,
    providesTags: result => (result ? [{type: TracetestApiTags.TEST_RUN, id: result?.id}] : []),
    transformResponse: (rawTestResult: TRawTestRun) => TestRun(rawTestResult),
    async onCacheEntryAdded(arg, {cacheDataLoaded, cacheEntryRemoved, updateCachedData}) {
      const listener: IListenerFunction<TRawTestRun> = data => {
        updateCachedData(() => TestRun(data.event));
      };
      await WebSocketService.initWebSocketSubscription({
        listener,
        resource: `test/${arg.testId}/run/${arg.runId}`,
        waitToCleanSubscription: cacheEntryRemoved,
        waitToInitSubscription: cacheDataLoaded,
      });
    },
  }),
  reRun: builder.mutation<TestRun, {testId: string; runId: string}>({
    query: ({testId, runId}) => ({
      url: `/tests/${testId}/run/${runId}/rerun`,
      method: HTTP_METHOD.POST,
    }),
    invalidatesTags: (response, error, {testId, runId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-LIST`},
      {type: TracetestApiTags.TEST_RUN, id: runId},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
    transformResponse: (rawTestRun: TRawTestRun) => TestRun(rawTestRun),
  }),
  dryRun: builder.mutation<AssertionResults, {testId: string; runId: string; testDefinition: Partial<TRawTestSpecs>}>({
    query: ({testId, runId, testDefinition}) => ({
      url: `/tests/${testId}/run/${runId}/dry-run`,
      method: HTTP_METHOD.PUT,
      body: testDefinition,
    }),
    transformResponse: (rawTestResults: TRawAssertionResults) => AssertionResults(rawTestResults),
  }),
  deleteRunById: builder.mutation<Test, {testId: string; runId: string}>({
    query: ({testId, runId}) => ({url: `/tests/${testId}/run/${runId}`, method: 'DELETE'}),
    invalidatesTags: (result, error, {testId}) => [
      {type: TracetestApiTags.TEST_RUN},
      {type: TracetestApiTags.TEST, id: `${testId}-LIST`},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  getJUnitByRunId: builder.query<string, {testId: string; runId: string}>({
    query: ({testId, runId}) => ({url: `/tests/${testId}/run/${runId}/junit.xml`, responseHandler: 'text'}),
    providesTags: (result, error, {testId, runId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-${runId}-junit`},
    ],
  }),
  getSelectedSpans: builder.query<SelectedSpans, {testId: string; runId: string; query: string}>({
    query: ({testId, runId, query}) => `/tests/${testId}/run/${runId}/select?query=${encodeURIComponent(query)}`,
    providesTags: (result, error, {query}) => (result ? [{type: TracetestApiTags.SPAN, id: `${query}-LIST`}] : []),
    transformResponse: (rawSpanList: TRawSelectedSpans) => SelectedSpans(rawSpanList),
  }),
});

export default TestRunEndpoint;
