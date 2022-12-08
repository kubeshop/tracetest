import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {
  SupportedDataStores,
  TDataStoreConfig,
  TRawDataStoreConfig,
  TTestConnectionRequest,
  TTestConnectionResponse,
} from 'types/Config.types';
// import Config from 'models/Config.model';
import DataStoreConfigMock from 'models/__mocks__/DataStoreConfig.mock';

const ConfigEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getDataStoreConfig: builder.query<TDataStoreConfig, unknown>({
    // query: () => '/config',
    query: () => '/tests',
    providesTags: () => [{type: TracetestApiTags.CONFIG, id: 'datastore'}],
    transformResponse: () =>
      DataStoreConfigMock.model({
        dataStores: [{name: 'jaeger', type: SupportedDataStores.JAEGER}],
        defaultDataStore: 'jaeger',
      }),
  }),
  updateDatastoreConfig: builder.mutation<undefined, TRawDataStoreConfig>({
    query: config => ({
      url: '/config/datastores',
      method: HTTP_METHOD.PUT,
      body: config,
    }),
    invalidatesTags: [{type: TracetestApiTags.CONFIG, id: 'datastore'}],
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
