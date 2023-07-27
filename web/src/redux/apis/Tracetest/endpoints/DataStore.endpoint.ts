import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import ConnectionResult from 'models/ConnectionResult.model';
import {TRawDataStore} from 'models/DataStore.model';
import DataStoreConfig from 'models/DataStoreConfig.model';
import {TConnectionResult, TRawConnectionResult, TTestConnectionRequest} from 'types/DataStore.types';
import TraceTestAPI from '../Tracetest.api';

const dataStoreEndpoints = TraceTestAPI.injectEndpoints({
  endpoints: builder => ({
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
    deleteDataStore: builder.mutation<undefined, void>({
      query: () => ({
        url: `/datastores/current`,
        method: HTTP_METHOD.DELETE,
        headers: {'content-type': 'application/json'},
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
  }),
});

export const {useGetDataStoreQuery, useUpdateDataStoreMutation, useDeleteDataStoreMutation, useTestConnectionMutation} =
  dataStoreEndpoints;
