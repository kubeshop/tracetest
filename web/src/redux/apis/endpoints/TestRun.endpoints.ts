import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import AssertionResults from 'models/AssertionResults.model';
import SelectedSpans from 'models/SelectedSpans.model';
import TestRun from 'models/TestRun.model';
import WebSocketService, {IListenerFunction} from 'services/WebSocket.service';
import {TAssertionResults, TRawAssertionResults} from 'types/Assertion.types';
import {TRawSelectedSpans, TSelectedSpans} from 'types/SelectedSpans.types';
import {TTest, TTestApiEndpointBuilder} from 'types/Test.types';
import {TRawTestRun, TTestRun} from 'types/TestRun.types';
import {TRawTestSpecs} from 'types/TestSpecs.types';

function getTotalCountFromHeaders(meta: any) {
  return Number(meta?.response?.headers.get('x-total-count') || 0);
}

const TestRunEndpoint = (builder: TTestApiEndpointBuilder) => ({
  runTest: builder.mutation<TTestRun, {testId: string; environmentId?: string}>({
    query: ({testId, environmentId}) => ({
      url: `/tests/${testId}/run`,
      method: HTTP_METHOD.POST,
      body: {
        environmentId,
      },
    }),
    invalidatesTags: (response, error, {testId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-LIST`},
      {type: TracetestApiTags.TEST, id: 'LIST'},
    ],
    transformResponse: (rawTestRun: TRawTestRun) => TestRun(rawTestRun),
  }),
  getRunList: builder.query<PaginationResponse<TTestRun>, {testId: string; take?: number; skip?: number}>({
    query: ({testId, take = 25, skip = 0}) => `/tests/${testId}/run?take=${take}&skip=${skip}`,
    providesTags: (result, error, {testId}) => [{type: TracetestApiTags.TEST_RUN, id: `${testId}-LIST`}],
    transformResponse: (rawTestResultList: TRawTestRun[], meta) => ({
      total: getTotalCountFromHeaders(meta),
      items: rawTestResultList.map(rawTestResult => TestRun(rawTestResult)),
    }),
  }),
  getRunById: builder.query<TTestRun, {runId: string; testId: string}>({
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
  reRun: builder.mutation<TTestRun, {testId: string; runId: string}>({
    query: ({testId, runId}) => ({
      url: `/tests/${testId}/run/${runId}/rerun`,
      method: HTTP_METHOD.POST,
    }),
    invalidatesTags: (response, error, {testId, runId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-LIST`},
      {type: TracetestApiTags.TEST_RUN, id: runId},
    ],
    transformResponse: (rawTestRun: TRawTestRun) => TestRun(rawTestRun),
  }),
  dryRun: builder.mutation<TAssertionResults, {testId: string; runId: string; testDefinition: Partial<TRawTestSpecs>}>({
    query: ({testId, runId, testDefinition}) => ({
      url: `/tests/${testId}/run/${runId}/dry-run`,
      method: HTTP_METHOD.PUT,
      body: testDefinition,
    }),
    transformResponse: (rawTestResults: TRawAssertionResults) => AssertionResults(rawTestResults),
  }),
  deleteRunById: builder.mutation<TTest, {testId: string; runId: string}>({
    query: ({testId, runId}) => ({url: `/tests/${testId}/run/${runId}`, method: 'DELETE'}),
    invalidatesTags: (result, error, {testId}) => [
      {type: TracetestApiTags.TEST_RUN},
      {type: TracetestApiTags.TEST, id: `${testId}-LIST`},
    ],
  }),
  getJUnitByRunId: builder.query<string, {testId: string; runId: string}>({
    query: ({testId, runId}) => ({url: `/tests/${testId}/run/${runId}/junit.xml`, responseHandler: 'text'}),
    providesTags: (result, error, {testId, runId}) => [
      {type: TracetestApiTags.TEST_RUN, id: `${testId}-${runId}-junit`},
    ],
  }),
  getSelectedSpans: builder.query<TSelectedSpans, {testId: string; runId: string; query: string}>({
    query: ({testId, runId, query}) => `/tests/${testId}/run/${runId}/select?query=${encodeURIComponent(query)}`,
    providesTags: (result, error, {query}) => (result ? [{type: TracetestApiTags.SPAN, id: `${query}-LIST`}] : []),
    transformResponse: (rawSpanList: TRawSelectedSpans) => SelectedSpans(rawSpanList),
  }),
});

export default TestRunEndpoint;
