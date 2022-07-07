import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {uniq} from 'lodash';
import WebSocketService, {IListenerFunction} from 'services/WebSocket.service';
import {TRawTest, TTest} from 'types/Test.types';
import {HTTP_METHOD} from '../../constants/Common.constants';
import AssertionResults from '../../models/AssertionResults.model';
import Test from '../../models/Test.model';
import TestDefinition from '../../models/TestDefinition.model';
import TestRun from '../../models/TestRun.model';
import {TAssertion, TAssertionResults, TRawAssertionResults} from '../../types/Assertion.types';
import {TDraftTest} from '../../types/Plugins.types';
import {TRawTestDefinition, TTestDefinition} from '../../types/TestDefinition.types';
import {TRawTestRun, TTestRun} from '../../types/TestRun.types';

const PATH = `${document.baseURI}api/`;

enum Tags {
  TEST = 'test',
  TEST_DEFINITION = 'testDefinition',
  TEST_RUN = 'testRun',
  SPAN = 'span',
}

const TraceTestAPI = createApi({
  reducerPath: 'tests',
  baseQuery: fetchBaseQuery({
    baseUrl: PATH,
  }),
  tagTypes: Object.values(Tags),
  endpoints: build => ({
    // Tests
    createTest: build.mutation<TTest, TDraftTest>({
      query: newTest => ({
        url: '/tests',
        method: HTTP_METHOD.POST,
        body: newTest,
      }),
      transformResponse: (rawTest: TRawTest) => Test(rawTest),
      invalidatesTags: [{type: Tags.TEST, id: 'LIST'}],
    }),
    editTest: build.mutation<TTest, {test: TDraftTest; testId: string}>({
      query: ({test, testId}) => ({
        url: `/tests/${testId}`,
        method: HTTP_METHOD.PUT,
        body: test,
      }),
      invalidatesTags: test => [
        {type: Tags.TEST, id: 'LIST'},
        {type: Tags.TEST, id: test?.id},
      ],
    }),
    getTestList: build.query<TTest[], void>({
      query: () => '/tests',
      providesTags: () => [{type: Tags.TEST, id: 'LIST'}],
      transformResponse: (rawTestList: TTest[]) => rawTestList.map(rawTest => Test(rawTest)),
    }),
    getTestById: build.query<TTest, {testId: string}>({
      query: ({testId}) => `/tests/${testId}`,
      providesTags: result => [{type: Tags.TEST, id: result?.id}],
      transformResponse: (rawTest: TRawTest) => Test(rawTest),
    }),
    deleteTestById: build.mutation<TTest, {testId: string}>({
      query: ({testId}) => ({url: `/tests/${testId}`, method: 'DELETE'}),
      invalidatesTags: [{type: Tags.TEST, id: 'LIST'}],
    }),

    // Test Definition
    getTestDefinition: build.query<TTestDefinition[], {testId: string}>({
      query: ({testId}) => `/tests/${testId}/definition`,
      providesTags: (result, error, {testId}) => [{type: Tags.TEST_DEFINITION, id: testId}],
      transformResponse: (testDefinitionList: TRawTestDefinition[]) =>
        testDefinitionList.map(rawTestDefinition => TestDefinition(rawTestDefinition)),
    }),
    setTestDefinition: build.mutation<TAssertion, {testId: string; testDefinition: Partial<TRawTestDefinition>}>({
      query: ({testId, testDefinition}) => ({
        url: `/tests/${testId}/definition`,
        method: HTTP_METHOD.PUT,
        body: testDefinition,
      }),
      invalidatesTags: (result, error, {testId}) => [
        {type: Tags.TEST, id: testId},
        {type: Tags.TEST_DEFINITION, id: testId},
      ],
    }),

    // Test Runs
    runTest: build.mutation<TTestRun, {testId: string}>({
      query: ({testId}) => ({
        url: `/tests/${testId}/run`,
        method: HTTP_METHOD.POST,
      }),
      invalidatesTags: (response, error, {testId}) => [
        {type: Tags.TEST_RUN, id: `${testId}-LIST`},
        {type: Tags.TEST, id: 'LIST'},
      ],
      transformResponse: (rawTestRun: TRawTestRun) => TestRun(rawTestRun),
    }),
    getRunList: build.query<TTestRun[], {testId: string; take?: number; skip?: number}>({
      query: ({testId, take = 25, skip = 0}) => `/tests/${testId}/run?take=${take}&skip=${skip}`,
      providesTags: (result, error, {testId}) => [{type: Tags.TEST_RUN, id: `${testId}-LIST`}],
      transformResponse: (rawTestResultList: TRawTestRun[]) =>
        rawTestResultList.map(rawTestResult => TestRun(rawTestResult)),
    }),
    getRunById: build.query<TTestRun, {runId: string; testId: string}>({
      query: ({testId, runId}) => `/tests/${testId}/run/${runId}`,
      providesTags: result => (result ? [{type: Tags.TEST_RUN, id: result?.id}] : []),
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
    reRun: build.mutation<TTestRun, {testId: string; runId: string}>({
      query: ({testId, runId}) => ({
        url: `/tests/${testId}/run/${runId}/rerun`,
        method: HTTP_METHOD.POST,
      }),
      invalidatesTags: (response, error, {testId, runId}) => [
        {type: Tags.TEST_RUN, id: `${testId}-LIST`},
        {type: Tags.TEST_RUN, id: runId},
      ],
      transformResponse: (rawTestRun: TRawTestRun) => TestRun(rawTestRun),
    }),
    dryRun: build.mutation<
      TAssertionResults,
      {testId: string; runId: string; testDefinition: Partial<TRawTestDefinition>}
    >({
      query: ({testId, runId, testDefinition}) => ({
        url: `/tests/${testId}/run/${runId}/dry-run`,
        method: HTTP_METHOD.PUT,
        body: testDefinition,
      }),
      transformResponse: (rawTestResults: TRawAssertionResults) => AssertionResults(rawTestResults),
    }),
    deleteRunById: build.mutation<TTest, {testId: string; runId: string}>({
      query: ({testId, runId}) => ({url: `/tests/${testId}/run/${runId}`, method: 'DELETE'}),
      invalidatesTags: (result, error, {testId}) => [{type: Tags.TEST_RUN, id: `${testId}-LIST`}],
    }),
    getJUnitByRunId: build.query<string, {testId: string; runId: string}>({
      query: ({testId, runId}) => ({url: `/tests/${testId}/run/${runId}/junit.xml`, responseHandler: 'text'}),
      providesTags: (result, error, {testId, runId}) => [{type: Tags.TEST_RUN, id: `${testId}-${runId}-junit`}],
    }),
    getTestDefinitionYamlByRunId: build.query<string, {testId: string; version: number}>({
      query: ({testId, version}) => ({
        url: `/tests/${testId}/version/${version}/definition.yaml`,
        responseHandler: 'text',
      }),
      providesTags: (result, error, {testId, version}) => [
        {type: Tags.TEST_RUN, id: `${testId}-${version}-definition`},
      ],
    }),

    // Spans
    getSelectedSpans: build.query<string[], {testId: string; runId: string; query: string}>({
      query: ({testId, runId, query}) => `/tests/${testId}/run/${runId}/select?query=${encodeURIComponent(query)}`,
      providesTags: (result, error, {query}) => (result ? [{type: Tags.SPAN, id: `${query}-LIST`}] : []),
      transformResponse: (rawSpanList: string[]) => uniq(rawSpanList),
    }),
  }),
});

export const {
  useCreateTestMutation,
  useGetTestByIdQuery,
  useGetTestListQuery,
  useRunTestMutation,
  useDeleteTestByIdMutation,
  useGetTestDefinitionQuery,
  useSetTestDefinitionMutation,
  useGetRunByIdQuery,
  useGetRunListQuery,
  useGetSelectedSpansQuery,
  useLazyGetSelectedSpansQuery,
  useReRunMutation,
  useLazyGetRunListQuery,
  useDryRunMutation,
  useDeleteRunByIdMutation,
  useGetJUnitByRunIdQuery,
  useLazyGetJUnitByRunIdQuery,
  useGetTestDefinitionYamlByRunIdQuery,
  useLazyGetTestDefinitionYamlByRunIdQuery,
  useEditTestMutation,
} = TraceTestAPI;
export const {endpoints} = TraceTestAPI;

export default TraceTestAPI;
