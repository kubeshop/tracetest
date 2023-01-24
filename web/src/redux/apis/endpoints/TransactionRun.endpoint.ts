import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import TransactionRun from 'models/TransactionRun.model';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {TRawTransactionRun, TTransactionRun} from 'types/TransactionRun.types';
import {IListenerFunction} from 'gateways/WebSocket.gateway';
import WebSocketService from 'services/WebSocket.service';
import {TEnvironmentValue} from 'types/Environment.types';
import {getTotalCountFromHeaders} from 'utils/Common';
import RunError from 'models/RunError.model';

const TransactionRunEndpoint = (builder: TTestApiEndpointBuilder) => ({
  runTransaction: builder.mutation<
    TTransactionRun,
    {transactionId: string; environmentId?: string; variables?: TEnvironmentValue[]}
  >({
    query: ({transactionId, environmentId, variables = []}) => ({
      url: `/transactions/${transactionId}/run`,
      method: HTTP_METHOD.POST,
      body: {environmentId, variables},
    }),
    invalidatesTags: (result, error, {transactionId}) => [
      {type: TracetestApiTags.TRANSACTION_RUN, id: `${transactionId}-LIST`},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
    transformResponse: (rawTransactionRun: TRawTransactionRun) => TransactionRun(rawTransactionRun),
    transformErrorResponse: ({data: result}) => RunError(result),
  }),

  getTransactionRuns: builder.query<
    PaginationResponse<TTransactionRun>,
    {transactionId: string; take?: number; skip?: number}
  >({
    query: ({transactionId, take = 25, skip = 0}) => `/transactions/${transactionId}/run?take=${take}&skip=${skip}`,
    providesTags: (result, error, {transactionId}) => [
      {type: TracetestApiTags.TRANSACTION_RUN, id: `${transactionId}-LIST`},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
    transformResponse: (rawTransactionRuns: TRawTransactionRun[], meta) => ({
      total: getTotalCountFromHeaders(meta),
      items: rawTransactionRuns.map(rawTransactionRun => TransactionRun(rawTransactionRun)),
    }),
  }),

  getTransactionRunById: builder.query<TTransactionRun, {transactionId: string; runId: string}>({
    query: ({transactionId, runId}) => `/transactions/${transactionId}/run/${runId}`,
    providesTags: result => [{type: TracetestApiTags.TRANSACTION_RUN, id: result?.id}],
    transformResponse: (rawTransactionRun: TRawTransactionRun) => TransactionRun(rawTransactionRun),
    async onCacheEntryAdded(arg, {cacheDataLoaded, cacheEntryRemoved, updateCachedData}) {
      const listener: IListenerFunction<TRawTransactionRun> = data => {
        updateCachedData(() => TransactionRun(data.event));
      };

      await WebSocketService.initWebSocketSubscription({
        listener,
        resource: `transaction/${arg.transactionId}/run/${arg.runId}`,
        waitToCleanSubscription: cacheEntryRemoved,
        waitToInitSubscription: cacheDataLoaded,
      });
    },
  }),

  deleteTransactionRunById: builder.mutation<TTransactionRun, {transactionId: string; runId: string}>({
    query: ({transactionId, runId}) => ({
      url: `/transactions/${transactionId}/run/${runId}`,
      method: HTTP_METHOD.DELETE,
    }),
    invalidatesTags: (result, error, {transactionId}) => [
      {type: TracetestApiTags.TRANSACTION_RUN, id: `${transactionId}-LIST`},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
});

export default TransactionRunEndpoint;
