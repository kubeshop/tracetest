import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {IAssertion, ITestAssertionResult} from 'types/Assertion.types';
import {ITest} from 'types/Test.types';
import {ITestRunResult, IRawTestRunResult} from 'types/TestRunResult.types';
import {TRecursivePartial} from 'types/Common.types';
import TestRunResult from 'models/TestRunResult.model';

export const testAPI = createApi({
  reducerPath: 'testsAPI',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/',
  }),
  tagTypes: ['Test', 'TestResult', 'Trace'],
  endpoints: build => ({
    createTest: build.mutation<ITest, TRecursivePartial<ITest>>({
      query: newTest => ({
        url: `/tests`,
        method: 'POST',
        body: newTest,
      }),
      invalidatesTags: [{type: 'Test', id: 'LIST'}],
    }),
    runTest: build.mutation<ITestRunResult, string>({
      query: testId => ({
        url: `/tests/${testId}/run`,
        method: 'POST',
      }),
      invalidatesTags: ['TestResult'],
    }),
    getTests: build.query<ITest[], void>({
      query: () => `/tests`,
      providesTags: result =>
        result
          ? [...result.map(i => ({type: 'Test' as const, id: i.testId})), {type: 'Test', id: 'LIST'}]
          : [{type: 'Test', id: 'LIST'}],
    }),
    getTestById: build.query<ITest, string>({
      query: id => `/tests/${id}`,
      providesTags: result => [{type: 'Test', id: result?.testId}],
    }),
    getTestAssertions: build.query<IAssertion[], string>({
      query: testId => `/tests/${testId}/assertions`,
    }),
    createAssertion: build.mutation<IAssertion, {testId: string; assertion: Partial<IAssertion>}>({
      query: ({testId, assertion}) => ({
        url: `/tests/${testId}/assertions`,
        method: 'POST',
        body: assertion,
      }),
      invalidatesTags: (result, error, args) => [{type: 'Test', id: args.testId}],
    }),
    updateAssertion: build.mutation<IAssertion, {testId: string; assertionId: string; assertion: Partial<IAssertion>}>({
      query: ({testId, assertion, assertionId}) => ({
        url: `/tests/${testId}/assertions/${assertionId}`,
        method: 'PUT',
        body: assertion,
      }),
      invalidatesTags: (result, error, args) => [{type: 'Test', id: args.testId}],
    }),
    getTestResults: build.query<ITestRunResult[], string>({
      query: id => `/tests/${id}/results`,
      providesTags: result =>
        result ? [{type: 'TestResult' as const, id: 'LIST'}] : [{type: 'TestResult' as const, id: 'LIST'}],
      transformResponse: (rawTestResultList: IRawTestRunResult[]) =>
        rawTestResultList.map(rawTestResult => TestRunResult(rawTestResult)),
    }),
    updateTestResult: build.mutation<
      ITestRunResult,
      {testId: string; resultId: string; assertionResult: ITestAssertionResult}
    >({
      query: ({testId, resultId, assertionResult}) => ({
        url: `/tests/${testId}/results/${resultId}`,
        method: 'PUT',
        body: assertionResult,
      }),
      invalidatesTags: (result, error, args) => [{type: 'TestResult', id: args.resultId}],
      transformResponse: (rawTestResult: IRawTestRunResult) => TestRunResult(rawTestResult),
    }),
    getTestResultById: build.query<ITestRunResult, Pick<ITest, 'testId'> & {resultId: string}>({
      query: ({testId, resultId}) => `/tests/${testId}/results/${resultId}`,
      providesTags: result => (result ? [{type: 'TestResult' as const, id: result?.resultId}] : []),
      transformResponse: (rawTestResult: IRawTestRunResult) => TestRunResult(rawTestResult),
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
