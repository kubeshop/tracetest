import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {TRecursivePartial} from 'types/Common.types';
import {TRawTest, TTest} from 'types/Test.types';
import {HTTP_METHOD} from '../../constants/Common.constants';
import Test from '../../models/Test.model';
import TestDefinition from '../../models/TestDefinition.model';
import TestRun from '../../models/TestRun.model';
import {TAssertion} from '../../types/Assertion.types';
import {TRawTestDefinition, TTestDefinition} from '../../types/TestDefinition.types';
import {TRawTestRun, TTestRun} from '../../types/TestRun.types';

const PATH = `${document.location.protocol}//${document.location.host}/api/`;

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
    createTest: build.mutation<TTest, TRecursivePartial<TTest>>({
      query: newTest => ({
        url: '/tests',
        method: HTTP_METHOD.POST,
        body: newTest,
      }),
      transformResponse: (rawTest: TRawTest) => Test(rawTest),
      invalidatesTags: [{type: Tags.TEST, id: 'LIST'}],
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

    // Spans
    getSelectedSpans: build.query<string[], {testId: string; runId: string; query: string}>({
      query: ({testId, runId, query}) => `/tests/${testId}/run/${runId}/select?query=${query}`,
      providesTags: (result, error, {query}) => (result ? [{type: Tags.SPAN, id: `${query}-LIST`}] : []),
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
  useReRunMutation,
  useLazyGetRunListQuery,
} = TraceTestAPI;
export const {endpoints} = TraceTestAPI;

export default TraceTestAPI;
