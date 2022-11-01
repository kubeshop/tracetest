import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {TracetestApiTags} from 'constants/Test.constants';
import {TTestApiEndpointBuilder} from '../../types/Test.types';
import EnvironmentEndpoint from './endpoints/Environment.endpoint';
import ExpressionEndpoint from './endpoints/Expression.endpoint';
import TestEndpoint from './endpoints/Test.endpoint';
import TestRunEndpoint from './endpoints/TestRun.endpoints';
import TransactionEndpoint from './endpoints/Transaction.endpoint';

const PATH = `${document.baseURI}api/`;

const TraceTestAPI = createApi({
  reducerPath: 'tests',
  baseQuery: fetchBaseQuery({
    baseUrl: PATH,
  }),
  tagTypes: Object.values(TracetestApiTags),
  endpoints(builder: TTestApiEndpointBuilder) {
    return {
      ...TransactionEndpoint(builder),
      ...TestRunEndpoint(builder),
      ...TestEndpoint(builder),
      ...EnvironmentEndpoint(builder),
      ...ExpressionEndpoint(builder),
    };
  },
});

export const {
  useCreateTestMutation,
  useGetTestByIdQuery,
  useGetTestVersionByIdQuery,
  useGetTestListQuery,
  useRunTestMutation,
  useDeleteTestByIdMutation,
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
  useGetEnvironmentsQuery,
  useCreateEnvironmentMutation,
  useUpdateEnvironmentMutation,
  useDeleteEnvironmentMutation,
  useGetTransactionByIdQuery,
  useGetTransactionRunByIdQuery,
  useDeleteTransactionByIdMutation,
  useParseExpressionMutation,
  useSetTestOutputsMutation,
} = TraceTestAPI;
export const {endpoints} = TraceTestAPI;

export default TraceTestAPI;
