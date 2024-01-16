import {BaseQueryFn, FetchArgs, FetchBaseQueryError, createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';
import {EndpointBuilder} from '@reduxjs/toolkit/dist/query/endpointDefinitions';
import {TracetestApiTags} from 'constants/Test.constants';
import Env from 'utils/Env';
import {dataStoreEndpoints} from './endpoints/DataStore.endpoint';
import {expressionEndpoints} from './endpoints/Expression.endpoint';
import {testEndpoints} from './endpoints/Test.endpoint';
import {testRunEndpoints} from './endpoints/TestRun.endpoint';
import {testSuiteEndpoints} from './endpoints/TestSuite.endpoint';
import {testSuiteRunEndpoints} from './endpoints/TestSuiteRun.endpoint';
import {resourceEndpoints} from './endpoints/Resource.endpoint';
import {variableSetEndpoints} from './endpoints/VariableSet.endpoint';
import {settingsEndpoints} from './endpoints/Setting.endpoint';
import {wizardEndpoints} from './endpoints/Wizard.endpoint';

export type TTestApiEndpointBuilder = EndpointBuilder<
  BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError>,
  TracetestApiTags,
  'tracetest'
>;

export type TBaseQueryFn = BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError>;

type TTracetestAPI = ReturnType<typeof generateAPI>;

const defaultBaseQuery: TBaseQueryFn = async (args, api, extraOptions) => {
  const baseUrl = Env.get('baseApiUrl');

  return fetchBaseQuery({
    baseUrl,
    prepareHeaders: headers => {
      headers.set('x-source', 'web');
    },
  })(args, api, extraOptions);
};

type TSingletonTracetestAPI = {
  instance: TTracetestAPI;
  create: (customBaseQuery?: typeof defaultBaseQuery) => void;
};

const generateAPI = (customBaseQuery = defaultBaseQuery) => {
  return createApi({
    reducerPath: 'tracetest',
    baseQuery: customBaseQuery,
    tagTypes: Object.values(TracetestApiTags),
    endpoints: builder => ({
      ...dataStoreEndpoints(builder),
      ...expressionEndpoints(builder),
      ...testEndpoints(builder),
      ...testRunEndpoints(builder),
      ...testSuiteEndpoints(builder),
      ...testSuiteRunEndpoints(builder),
      ...resourceEndpoints(builder),
      ...variableSetEndpoints(builder),
      ...settingsEndpoints(builder),
      ...wizardEndpoints(builder),
    }),
  });
};

const TracetestAPI: TSingletonTracetestAPI = {
  instance: {endpoints: {}} as TTracetestAPI,
  create(customBaseQuery = defaultBaseQuery) {
    this.instance = generateAPI(customBaseQuery);
  },
};

export default TracetestAPI;
