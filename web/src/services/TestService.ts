import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {Assertion, Test} from '../types';

export const testAPI = createApi({
  reducerPath: 'testsAPI',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/',
  }),
  tagTypes: ['Test'],
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
    }),
    getTests: build.query<Test[], void>({
      query: () => `/tests`,
      providesTags: [{type: 'Test', id: 'LIST'}],
    }),
    getTestById: build.query<Test, string>({
      query: id => `/tests/${id}`,
    }),
    getTestAssertions: build.query<Assertion[], Pick<Test, 'id'>>({
      query: ({id}) => `/tests/${id}/assertions`,
    }),
    getTestResults: build.query<Assertion[], Pick<Test, 'id'>>({
      query: ({id}) => `/tests/${id}/results`,
    }),
    getTestResultById: build.query<Assertion[], Pick<Test, 'id'> & {resultId: string}>({
      query: ({id, resultId}) => `/tests/${id}/results/${resultId}`,
    }),
  }),
});

export const {
  useCreateTestMutation,
  useGetTestAssertionsQuery,
  useGetTestByIdQuery,
  useGetTestResultByIdQuery,
  useGetTestResultsQuery,
  useGetTestsQuery,
} = testAPI;
