import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {ITest} from 'types/Test.types';
import {IRawTestRunResult, ITestRunResult} from 'types/TestRunResult.types';
import {TRecursivePartial} from 'types/Common.types';
import {IAssertion, ITestAssertionResult} from '../../types/Assertion.types';
import TestRunResult from '../../models/TestRunResult.model';

const TestAPI = createApi({
  reducerPath: 'tests',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/',
  }),
  tagTypes: ['Test', 'Assertion', 'TestRunResult'],
  endpoints: build => ({
    // Tests
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
    }),
    getTestList: build.query<ITest[], void>({
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
    deleteTestById: build.mutation<ITest, string>({
      query: id => ({url: `/tests/${id}`, method: 'DELETE'}),
      invalidatesTags: [{type: 'Test', id: 'LIST'}],
    }),
    // Assertions
    getAssertions: build.query<IAssertion[], string>({
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

    // Test Results
    getResultList: build.query<ITestRunResult[], string>({
      query: id => `/tests/${id}/results`,
      providesTags: result =>
        result ? [{type: 'TestRunResult' as const, id: 'LIST'}] : [{type: 'TestRunResult' as const, id: 'LIST'}],
      transformResponse: (rawTestResultList: IRawTestRunResult[]) =>
        rawTestResultList.map(rawTestResult => TestRunResult(rawTestResult)),
    }),
    updateResult: build.mutation<
      ITestRunResult,
      {testId: string; resultId: string; assertionResult: ITestAssertionResult}
    >({
      query: ({testId, resultId, assertionResult}) => ({
        url: `/tests/${testId}/results/${resultId}`,
        method: 'PUT',
        body: assertionResult,
      }),
      invalidatesTags: (result, error, args) => [{type: 'TestRunResult', id: args.resultId}],
      transformResponse: (rawTestResult: IRawTestRunResult) => TestRunResult(rawTestResult),
    }),
    getResultById: build.query<ITestRunResult, Pick<ITest, 'testId'> & {resultId: string}>({
      query: ({testId, resultId}) => `/tests/${testId}/results/${resultId}`,
      providesTags: result => (result ? [{type: 'TestRunResult' as const, id: result?.resultId}] : []),
      transformResponse: (rawTestResult: IRawTestRunResult) => TestRunResult(rawTestResult),
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
  useDeleteTestByIdMutation,
  useUpdateResultMutation,
} = TestAPI;
export const {endpoints} = TestAPI;

export default TestAPI;
