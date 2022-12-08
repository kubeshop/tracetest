import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {
  TConfig,
  TTestConnectionRequest,
  TTestConnectionResponse,
  TUpdateDataStoreConfigRequest,
} from 'types/Config.types';
// import Config from 'models/Config.model';
import ConfigMock from 'models/__mocks__/Config.mock';

const ConfigEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getConfig: builder.query<TConfig, unknown>({
    // query: () => '/config',
    query: () => '/tests',
    providesTags: () => [{type: TracetestApiTags.CONFIG, id: 'config'}],
    transformResponse: () =>
      ConfigMock.model({
        telemetry: {dataStores: [{type: 'jaeger'}]},
        server: {telemetry: {dataStore: 'jaeger'}},
      }),
  }),
  updateDatastoreConfig: builder.mutation<undefined, TUpdateDataStoreConfigRequest>({
    query: config => ({
      url: '/config/datastore',
      method: HTTP_METHOD.PUT,
      body: config,
    }),
    invalidatesTags: [{type: TracetestApiTags.CONFIG, id: 'config'}],
  }),
  testConnection: builder.mutation<TTestConnectionResponse, TTestConnectionRequest>({
    query: connectionTest => ({
      url: `/config/connection`,
      method: HTTP_METHOD.PUT,
      body: connectionTest,
    }),
  }),
});

export default ConfigEndpoint;
