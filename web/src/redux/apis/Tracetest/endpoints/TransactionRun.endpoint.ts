import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import {TEnvironmentValue} from 'models/Environment.model';
import RunError from 'models/RunError.model';
import TransactionRun, {TRawTransactionResourceRun} from 'models/TransactionRun.model';
import WebSocketService, {IListenerFunction} from 'services/WebSocket.service';
import {getTotalCountFromHeaders} from 'utils/Common';
import TraceTestAPI from '../Tracetest.api';

const transactionRunEndpoints = TraceTestAPI.injectEndpoints({
  endpoints: builder => ({
    runTransaction: builder.mutation<
      TransactionRun,
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
      transformResponse: (rawTransactionRun: TRawTransactionResourceRun) => TransactionRun(rawTransactionRun),
      transformErrorResponse: ({data: result}) => RunError(result),
    }),

    getTransactionRuns: builder.query<
      PaginationResponse<TransactionRun>,
      {transactionId: string; take?: number; skip?: number}
    >({
      query: ({transactionId, take = 25, skip = 0}) => `/transactions/${transactionId}/run?take=${take}&skip=${skip}`,
      providesTags: (result, error, {transactionId}) => [
        {type: TracetestApiTags.TRANSACTION_RUN, id: `${transactionId}-LIST`},
        {type: TracetestApiTags.RESOURCE, id: 'LIST'},
      ],
      transformResponse: (rawTransactionRuns: TRawTransactionResourceRun[], meta) => ({
        total: getTotalCountFromHeaders(meta),
        items: rawTransactionRuns.map(rawTransactionRun => TransactionRun(rawTransactionRun)),
      }),
    }),

    getTransactionRunById: builder.query<TransactionRun, {transactionId: string; runId: string}>({
      query: ({transactionId, runId}) => `/transactions/${transactionId}/run/${runId}`,
      providesTags: result => [{type: TracetestApiTags.TRANSACTION_RUN, id: result?.id}],
      transformResponse: (rawTransactionRun: TRawTransactionResourceRun) => TransactionRun(rawTransactionRun),
      async onCacheEntryAdded(arg, {cacheDataLoaded, cacheEntryRemoved, updateCachedData}) {
        const listener: IListenerFunction<TRawTransactionResourceRun> = data => {
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

    deleteTransactionRunById: builder.mutation<TransactionRun, {transactionId: string; runId: string}>({
      query: ({transactionId, runId}) => ({
        url: `/transactions/${transactionId}/run/${runId}`,
        method: HTTP_METHOD.DELETE,
      }),
      invalidatesTags: (result, error, {transactionId}) => [
        {type: TracetestApiTags.TRANSACTION_RUN, id: `${transactionId}-LIST`},
        {type: TracetestApiTags.RESOURCE, id: 'LIST'},
      ],
    }),
  }),
});

export const {
  useGetTransactionRunsQuery,
  useLazyGetTransactionRunsQuery,
  useGetTransactionRunByIdQuery,
  useDeleteTransactionRunByIdMutation,
  useRunTransactionMutation,
} = transactionRunEndpoints;
