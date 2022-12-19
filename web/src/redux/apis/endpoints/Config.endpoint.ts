import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {
  SupportedDataStores,
  TDataStoreConfig,
  TRawDataStoreConfig,
  TTestConnectionRequest,
  TConnectionResult,
  TRawConnectionResult,
} from 'types/Config.types';
import DataStoreConfigMock from 'models/__mocks__/DataStoreConfig.mock';
import ConnectionResult from '../../../models/ConnectionResult.model';

const ConfigEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getDataStoreConfig: builder.query<TDataStoreConfig, unknown>({
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
  testConnection: builder.mutation<TConnectionResult, TTestConnectionRequest>({
    query: connectionTest => ({
      url: `/config/connection`,
      method: HTTP_METHOD.POST,
      body: connectionTest,
    }),
    transformResponse: (result: TRawConnectionResult) => ConnectionResult(result),
    transformErrorResponse: ({data: result}) => ConnectionResult(result as TRawConnectionResult),
  }),
});

export default ConfigEndpoint;
