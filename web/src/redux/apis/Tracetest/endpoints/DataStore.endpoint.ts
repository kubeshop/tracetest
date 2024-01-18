import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import ConnectionResult from 'models/ConnectionResult.model';
import {TRawDataStore} from 'models/DataStore.model';
import DataStoreConfig from 'models/DataStoreConfig.model';
import OTLPTestConnectionResponse, {TRawOTLPTestConnectionResponse} from 'models/OTLPTestConnectionResponse.model';
import {TConnectionResult, TRawConnectionResult, TTestConnectionRequest} from 'types/DataStore.types';
import {TTestApiEndpointBuilder} from '../Tracetest.api';

export const dataStoreEndpoints = (builder: TTestApiEndpointBuilder) => ({
  getDataStore: builder.query<DataStoreConfig, unknown>({
    query: () => ({
      url: '/datastores/current',
      method: HTTP_METHOD.GET,
      headers: {'content-type': 'application/json'},
    }),
    providesTags: () => [{type: TracetestApiTags.DATA_STORE, id: 'datastore'}],
    transformResponse: (rawDataStore: TRawDataStore) => DataStoreConfig(rawDataStore),
  }),
  updateDataStore: builder.mutation<undefined, {dataStore: TRawDataStore}>({
    query: ({dataStore}) => ({
      url: `/datastores/current`,
      method: HTTP_METHOD.PUT,
      body: dataStore,
    }),
    invalidatesTags: [{type: TracetestApiTags.DATA_STORE, id: 'datastore'}],
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

  testOtlpConnection: builder.query<OTLPTestConnectionResponse, unknown>({
    query: () => ({
      url: '/config/connection/otlp',
      method: HTTP_METHOD.GET,
      headers: {'content-type': 'application/json'},
    }),
    providesTags: () => [{type: TracetestApiTags.DATA_STORE, id: 'datastore'}],
    transformResponse: (raw: TRawOTLPTestConnectionResponse) => OTLPTestConnectionResponse(raw),
  }),
  resetTestOtlpConnection: builder.mutation<undefined, undefined>({
    query: () => ({
      url: `/config/connection/otlp/reset`,
      method: HTTP_METHOD.POST,
    }),
  }),
});
