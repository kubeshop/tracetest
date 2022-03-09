import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {Assertion, ITestResult, ITrace, Test, TestId} from 'types';

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
          ? [...result.map(i => ({type: 'Test' as const, id: i.id})), {type: 'Test', id: 'LIST'}]
          : [{type: 'Test', id: 'LIST'}],
    }),
    getTestById: build.query<Test, string>({
      query: id => `/tests/${id}`,
      providesTags: result => [{type: 'Test', id: result?.id}],
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
          ? [...result.map(el => ({type: 'TestResult' as const, id: el.id})), {type: 'TestResult' as const, id: 'LIST'}]
          : [{type: 'TestResult' as const, id: 'LIST'}],
    }),
    getTestResultById: build.query<Assertion[], Pick<Test, 'id'> & {resultId: string}>({
      query: ({id, resultId}) => `/tests/${id}/results/${resultId}`,
    }),
    getTestTrace: build.query<ITrace, Pick<Test, 'id'> & {resultId: string}>({
      query: ({id, resultId}) => `/tests/${id}/results/${resultId}/trace`,
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
  useGetTestTraceQuery,
} = testAPI;
