import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {TRecursivePartial} from 'types/Common.types';
import {ITest} from 'types/Test.types';
import {IRawTestRunResult, ITestRunResult} from 'types/TestRunResult.types';
import {HTTP_METHOD} from '../../constants/Common.constants';
import TestRunResult from '../../models/TestRunResult.model';
import {IAssertion, ITestAssertionResult} from '../../types/Assertion.types';

const PATH = `${document.location.protocol}//${document.location.host}/api/`;

const TestAPI = createApi({
  reducerPath: 'tests',
  baseQuery: fetchBaseQuery({
    baseUrl: PATH,
  }),
  tagTypes: ['Test', 'Assertion', 'TestRunResult'],
  endpoints: build => ({
    // Tests
    createTest: build.mutation<ITest, TRecursivePartial<ITest>>({
      query: newTest => ({
        url: `/tests`,
        method: HTTP_METHOD.POST,
        body: newTest,
      }),
      invalidatesTags: [{type: 'Test', id: 'LIST'}],
    }),
    runTest: build.mutation<ITestRunResult, string>({
      query: testId => ({
        url: `/tests/${testId}/run`,
        method: HTTP_METHOD.POST,
      }),
      invalidatesTags: (response, error, testId) => [
        {type: 'TestRunResult', id: `${testId}-LIST`},
        {type: 'Test', id: 'LIST'},
      ],
    }),
    getTestList: build.query<ITest[], void>({
      query: () => `/tests`,
      providesTags: () => [{type: 'Test', id: 'LIST'}],
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
        method: HTTP_METHOD.POST,
        body: assertion,
      }),
      invalidatesTags: (result, error, args) => [{type: 'Test', id: args.testId}],
    }),
    updateAssertion: build.mutation<IAssertion, {testId: string; assertionId: string; assertion: Partial<IAssertion>}>({
      query: ({testId, assertion, assertionId}) => ({
        url: `/tests/${testId}/assertions/${assertionId}`,
        method: HTTP_METHOD.PUT,
        body: assertion,
      }),
      invalidatesTags: (result, error, args) => [{type: 'Test', id: args.testId}],
    }),
    deleteAssertion: build.mutation<IAssertion, {testId: string; assertionId: string}>({
      query: ({testId, assertionId}) => ({
        url: `/tests/${testId}/assertions/${assertionId}`,
        method: HTTP_METHOD.DELETE,
      }),
      invalidatesTags: (result, error, args) => [{type: 'Test', id: args.testId}],
    }),

    // Test Results
    getResultList: build.query<ITestRunResult[], {testId: string; take?: number; skip?: number}>({
      query: ({testId, take = 25, skip = 0}) => `/tests/${testId}/results?take=${take}&skip=${skip}`,
      providesTags: (result, error, {testId}) => [{type: 'TestRunResult' as const, id: `${testId}-LIST`}],
      transformResponse: (rawTestResultList: IRawTestRunResult[]) =>
        rawTestResultList.map(rawTestResult => TestRunResult(rawTestResult)),
    }),
    updateResult: build.mutation<
      ITestRunResult,
      {testId: string; resultId: string; assertionResult: ITestAssertionResult}
    >({
      query: ({testId, resultId, assertionResult}) => ({
        url: `/tests/${testId}/results/${resultId}`,
        method: HTTP_METHOD.PUT,
        body: assertionResult,
      }),
      invalidatesTags: (result, error, {testId, resultId}) => [
        {type: 'TestRunResult', id: resultId},
        {type: 'TestRunResult', id: `${testId}-LIST`},
      ],
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
  useGetResultByIdQuery,
  useGetResultListQuery,
  useGetTestByIdQuery,
  useGetTestListQuery,
  useRunTestMutation,
  useUpdateAssertionMutation,
  useDeleteTestByIdMutation,
  useGetAssertionsQuery,
  useUpdateResultMutation,
  useLazyGetResultListQuery,
  useDeleteAssertionMutation,
} = TestAPI;
export const {endpoints} = TestAPI;

export default TestAPI;
