import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {TAssertion, TTestAssertionResult} from '../entities/Assertion/Assertion.types';
import {TTest} from '../entities/Test/Test.types';
import {TTestRunResult} from '../entities/TestRunResult/TestRunResult.types';
import {TRecursivePartial} from '../types/Common.types';

export const testAPI = createApi({
  reducerPath: 'testsAPI',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/',
  }),
  tagTypes: ['Test', 'TestResult', 'Trace'],
  endpoints: build => ({
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
      invalidatesTags: ['TestResult'],
    }),
    getTests: build.query<TTest[], void>({
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
    getTestAssertions: build.query<TAssertion[], string>({
      query: testId => `/tests/${testId}/assertions`,
    }),
    createAssertion: build.mutation<TAssertion, {testId: string; assertion: Partial<TAssertion>}>({
      query: ({testId, assertion}) => ({
        url: `/tests/${testId}/assertions`,
        method: 'POST',
        body: assertion,
      }),
      invalidatesTags: (result, error, args) => [{type: 'Test', id: args.testId}],
    }),
    updateAssertion: build.mutation<TAssertion, {testId: string; assertionId: string; assertion: Partial<TAssertion>}>({
      query: ({testId, assertion, assertionId}) => ({
        url: `/tests/${testId}/assertions/${assertionId}`,
        method: 'PUT',
        body: assertion,
      }),
      invalidatesTags: (result, error, args) => [{type: 'Test', id: args.testId}],
    }),
    getTestResults: build.query<TTestRunResult[], string>({
      query: id => `/tests/${id}/results`,
      providesTags: result =>
        result ? [{type: 'TestResult' as const, id: 'LIST'}] : [{type: 'TestResult' as const, id: 'LIST'}],
    }),
    updateTestResult: build.mutation<
      TTestRunResult,
      {testId: string; resultId: string; assertionResult: TTestAssertionResult}
    >({
      query: ({testId, resultId, assertionResult}) => ({
        url: `/tests/${testId}/results/${resultId}`,
        method: 'PUT',
        body: assertionResult,
      }),
      invalidatesTags: (result, error, args) => [{type: 'TestResult', id: args.resultId}],
    }),
    getTestResultById: build.query<TTestRunResult, Pick<TTest, 'testId'> & {resultId: string}>({
      query: ({testId, resultId}) => `/tests/${testId}/results/${resultId}`,
      providesTags: result => (result ? [{type: 'TestResult' as const, id: result?.resultId}] : []),
    }),
  }),
});

export const {
  useCreateTestMutation,
  useCreateAssertionMutation,
  useUpdateTestResultMutation,
  useRunTestMutation,
  useGetTestAssertionsQuery,
  useGetTestByIdQuery,
  useGetTestResultByIdQuery,
  useGetTestResultsQuery,
  useGetTestsQuery,
  useUpdateAssertionMutation,
} = testAPI;
