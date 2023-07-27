import {endpoints as environmentEndpoints} from './endpoints/Environment.endpoint';
import {endpoints as testEndpoints} from './endpoints/Test.endpoint';
import {endpoints as testRunEndpoints} from './endpoints/TestRun.endpoint';

// eslint-disable-next-line no-restricted-exports
export {default} from './Tracetest.api';

export {
  useGetDataStoreQuery,
  useUpdateDataStoreMutation,
  useDeleteDataStoreMutation,
  useTestConnectionMutation,
} from './endpoints/DataStore.endpoint';

export {
  useGetEnvironmentsQuery,
  useCreateEnvironmentMutation,
  useUpdateEnvironmentMutation,
  useDeleteEnvironmentMutation,
} from './endpoints/Environment.endpoint';

export {useParseExpressionMutation} from './endpoints/Expression.endpoint';

export {
  useGetResourcesQuery,
  useGetResourceDefinitionQuery,
  useLazyGetResourceDefinitionQuery,
} from './endpoints/Resource.endpoint';

export {
  useGetConfigQuery,
  useGetPollingQuery,
  useGetDemoQuery,
  useGetLinterQuery,
  useGetTestRunnerQuery,
  useCreateSettingMutation,
  useUpdateSettingMutation,
} from './endpoints/Setting.endpoint';

export {
  useCreateTestMutation,
  useGetTestByIdQuery,
  useGetTestVersionByIdQuery,
  useGetTestListQuery,
  useDeleteTestByIdMutation,
  useEditTestMutation,
} from './endpoints/Test.endpoint';

export {
  useGetRunByIdQuery,
  useGetRunEventsQuery,
  useGetRunListQuery,
  useGetSelectedSpansQuery,
  useLazyGetSelectedSpansQuery,
  useRunTestMutation,
  useReRunMutation,
  useLazyGetRunListQuery,
  useDryRunMutation,
  useDeleteRunByIdMutation,
  useStopRunMutation,
  useGetJUnitByRunIdQuery,
  useLazyGetJUnitByRunIdQuery,
} from './endpoints/TestRun.endpoint';

export {
  useCreateTransactionMutation,
  useGetTransactionByIdQuery,
  useDeleteTransactionByIdMutation,
  useEditTransactionMutation,
  useGetTransactionVersionByIdQuery,
  useLazyGetTransactionVersionByIdQuery,
} from './endpoints/Transaction.endpoint';

export {
  useGetTransactionRunsQuery,
  useLazyGetTransactionRunsQuery,
  useGetTransactionRunByIdQuery,
  useDeleteTransactionRunByIdMutation,
  useRunTransactionMutation,
} from './endpoints/TransactionRun.endpoint';

export const endpoints = {...environmentEndpoints, ...testEndpoints, ...testRunEndpoints};
