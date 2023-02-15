import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import ConnectionResult from 'models/ConnectionResult.model';
import DataStoreConfig from 'models/DataStoreConfig.model';
import {TConnectionResult, TRawConnectionResult, TTestConnectionRequest} from 'types/DataStore.types';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {TRawDataStore} from 'models/DataStore.model';

const ConfigEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getDataStores: builder.query<DataStoreConfig, unknown>({
    query: () => '/datastores?take=50',
    providesTags: () => [{type: TracetestApiTags.CONFIG, id: 'datastore'}],
    transformResponse: (rawDataStores: TRawDataStore[]) => DataStoreConfig(rawDataStores),
  }),
  createDataStore: builder.mutation<undefined, TRawDataStore>({
    query: dataStore => ({
      url: '/datastores',
      method: HTTP_METHOD.POST,
      body: dataStore,
    }),
    invalidatesTags: [{type: TracetestApiTags.CONFIG, id: 'datastore'}],
  }),
  updateDataStore: builder.mutation<undefined, {dataStore: TRawDataStore; dataStoreId: string}>({
    query: ({dataStore, dataStoreId}) => ({
      url: `/datastores/${dataStoreId}`,
      method: HTTP_METHOD.PUT,
      body: dataStore,
    }),
    invalidatesTags: [{type: TracetestApiTags.CONFIG, id: 'datastore'}],
  }),
  deleteDataStore: builder.mutation<undefined, {dataStoreId: string}>({
    query: ({dataStoreId}) => ({
      url: `/datastores/${dataStoreId}`,
      method: HTTP_METHOD.DELETE,
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
