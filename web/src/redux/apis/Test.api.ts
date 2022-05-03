import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {TTest} from 'types/Test.types';
import {TRawTestRunResult, TTestRunResult} from 'types/TestRunResult.types';
import {TRecursivePartial} from 'types/Common.types';
import {TAssertion, TTestAssertionResult} from '../../types/Assertion.types';
import TestRunResult from '../../models/TestRunResult.model';

const TestAPI = createApi({
  reducerPath: 'tests',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/',
  }),
  tagTypes: ['Test', 'Assertion', 'TestRunResult'],
  endpoints: build => ({
    // Tests
    createTest: build.mutation<TTest, TRecursivePartial<TTest>>({
      query: newTest => ({
        url: `/tests`,
        method: 'POST',
        body: newTest,
      }),
      invalidatesTags: [{type: 'Test', id: 'LIST'}],
    }),
    runTest: build.mutation<TTestRunResult, string>({
      query: testId => ({
        url: `/tests/${testId}/run`,
        method: 'POST',
      }),
    }),
    getTestList: build.query<TTest[], void>({
      query: () => `/tests`,
      providesTags: result =>
        result
          ? [...result.map(i => ({type: 'Test' as const, id: i.testId})), {type: 'Test', id: 'LIST'}]
          : [{type: 'Test', id: 'LIST'}],
    }),
    getTestById: build.query<TTest, string>({
      query: id => `/tests/${id}`,
      providesTags: result => [{type: 'Test', id: result?.testId}],
    }),

    // Assertions
    getAssertions: build.query<TAssertion[], string>({
      query: testId => `/tests/${testId}/assertions`,
    }),
    createAssertion: build.mutation<TAssertion, {testId: string; assertion: TRecursivePartial<TAssertion>}>({
      query: ({testId, assertion}) => ({
        url: `/tests/${testId}/assertions`,
        method: 'POST',
        body: assertion,
      }),
      invalidatesTags: (result, error, args) => [{type: 'Test', id: args.testId}],
    }),
    updateAssertion: build.mutation<TAssertion, {testId: string; assertionId: string; assertion: TRecursivePartial<TAssertion>}>({
      query: ({testId, assertion, assertionId}) => ({
        url: `/tests/${testId}/assertions/${assertionId}`,
        method: 'PUT',
        body: assertion,
      }),
      invalidatesTags: (result, error, args) => [{type: 'Test', id: args.testId}],
    }),

    // Test Results
    getResultList: build.query<TTestRunResult[], string>({
      query: id => `/tests/${id}/results`,
      providesTags: result =>
        result ? [{type: 'TestRunResult' as const, id: 'LIST'}] : [{type: 'TestRunResult' as const, id: 'LIST'}],
      transformResponse: (rawTestResultList: TRawTestRunResult[]) =>
        rawTestResultList.map(rawTestResult => TestRunResult(rawTestResult)),
    }),
    updateResult: build.mutation<
      TTestRunResult,
      {testId: string; resultId: string; assertionResult: TTestAssertionResult}
    >({
      query: ({testId, resultId, assertionResult}) => ({
        url: `/tests/${testId}/results/${resultId}`,
        method: 'PUT',
        body: assertionResult,
      }),
      invalidatesTags: (result, error, args) => [{type: 'TestRunResult', id: args.resultId}],
      transformResponse: (rawTestResult: TRawTestRunResult) => TestRunResult(rawTestResult),
    }),
    getResultById: build.query<TTestRunResult, Pick<TTest, 'testId'> & {resultId: string}>({
      query: ({testId, resultId}) => `/tests/${testId}/results/${resultId}`,
      providesTags: result => (result ? [{type: 'TestRunResult' as const, id: result?.resultId}] : []),
      transformResponse: (rawTestResult: TRawTestRunResult) => TestRunResult(rawTestResult),
    }),
  }),
});

export const {
  useCreateAssertionMutation,
  useCreateTestMutation,
  useGetAssertionsQuery,
  useGetResultByIdQuery,
  useGetResultListQuery,
  useGetTestByIdQuery,
  useGetTestListQuery,
  useRunTestMutation,
  useUpdateAssertionMutation,
  useUpdateResultMutation,
} = TestAPI;
export const {endpoints} = TestAPI;

export default TestAPI;
