import {endpoints as variableSetEndpoints} from './endpoints/VariableSet.endpoint';
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
  useCreateVariableSetMutation,
  useDeleteVariableSetMutation,
  useGetVariableSetsQuery,
  useUpdateVariableSetMutation,
} from './endpoints/VariableSet.endpoint';

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
  useCreateTestSuiteMutation,
  useGetTestSuiteByIdQuery,
  useGetTestSuiteVersionByIdQuery,
  useLazyGetTestSuiteByIdQuery,
  useDeleteTestSuiteByIdMutation,
  useEditTestSuiteMutation,
  useLazyGetTestSuiteVersionByIdQuery,
} from './endpoints/TestSuite.endpoint';

export {
  useDeleteTestSuiteRunByIdMutation,
  useGetTestSuiteRunByIdQuery,
  useGetTestSuiteRunsQuery,
  useRunTestSuiteMutation,
  useLazyGetTestSuiteRunsQuery,
  useLazyGetTestSuiteRunByIdQuery,
} from './endpoints/TestSuiteRun.endpoint';

export const endpoints = {...variableSetEndpoints, ...testEndpoints, ...testRunEndpoints};
