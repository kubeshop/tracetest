import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {Assertion, ITestResult, Test, TestId, TestRunResult} from 'types';

export const testAPI = createApi({
  reducerPath: 'testsAPI',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/',
  }),
  tagTypes: ['Test', 'TestResult', 'Trace'],
  endpoints: build => ({
    createTest: build.mutation<Test, Partial<Test>>({
      query: newTest => ({
        url: `/tests`,
        method: 'POST',
        body: newTest,
      }),
      invalidatesTags: [{type: 'Test', id: 'LIST'}],
    }),
    runTest: build.mutation<{id: string}, string>({
      query: testId => ({
        url: `/tests/${testId}/run`,
        method: 'POST',
      }),
      invalidatesTags: ['TestResult'],
    }),
    getTests: build.query<Test[], void>({
      query: () => `/tests`,
      providesTags: result =>
        result
          ? [...result.map(i => ({type: 'Test' as const, id: i.testId})), {type: 'Test', id: 'LIST'}]
          : [{type: 'Test', id: 'LIST'}],
    }),
    getTestById: build.query<Test, string>({
      query: id => `/tests/${id}`,
      providesTags: result => [{type: 'Test', id: result?.testId}],
    }),
    getTestAssertions: build.query<Assertion[], string>({
      query: testId => `/tests/${testId}/assertions`,
    }),
    createAssertion: build.mutation<Assertion, {testId: string} & Partial<Assertion>>({
      query: ({testId, ...assertion}) => ({
        url: `/tests/${testId}/assertions`,
        method: 'POST',
        body: assertion,
      }),
      invalidatesTags: (result, error, args) => [{type: 'Test', id: args.testId}],
    }),
    getTestResults: build.query<ITestResult[], TestId>({
      query: id => `/tests/${id}/results`,
      providesTags: result =>
        result
          ? [
              ...result.map(el => ({type: 'TestResult' as const, id: el.resultId})),
              {type: 'TestResult' as const, id: 'LIST'},
            ]
          : [{type: 'TestResult' as const, id: 'LIST'}],
    }),
    getTestResultById: build.query<TestRunResult[], Pick<Test, 'testId'> & {resultId: string}>({
      query: ({testId, resultId}) => `/tests/${testId}/results/${resultId}`,
    }),
  }),
});

export const {
  useCreateTestMutation,
  useCreateAssertionMutation,
  useRunTestMutation,
  useGetTestAssertionsQuery,
  useGetTestByIdQuery,
  useGetTestResultByIdQuery,
  useGetTestResultsQuery,
  useGetTestsQuery,
} = testAPI;
