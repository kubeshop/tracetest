import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import TransactionRun from 'models/TransactionRun.model';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {TRawTransactionRun, TTransactionRun} from 'types/TransactionRun.types';
import {getTotalCountFromHeaders} from '../../../utils/Common';

const TransactionRunEndpoint = (builder: TTestApiEndpointBuilder) => ({
  runTransaction: builder.mutation<TTransactionRun, {transactionId: string; environmentId?: string}>({
    query: ({transactionId, environmentId}) => ({
      url: `/transactions/${transactionId}/run`,
      method: HTTP_METHOD.POST,
      body: {environmentId},
    }),
    invalidatesTags: (result, error, {transactionId}) => [
      {type: TracetestApiTags.TRANSACTION_RUN, id: `${transactionId}-LIST`},
    ],
    transformResponse: (rawTransactionRun: TRawTransactionRun) => TransactionRun(rawTransactionRun),
  }),

  getTransactionRuns: builder.query<
    PaginationResponse<TTransactionRun>,
    {transactionId: string; take?: number; skip?: number}
  >({
    query: ({transactionId, take = 25, skip = 0}) => `/transactions/${transactionId}/run?take=${take}&skip=${skip}`,
    providesTags: (result, error, {transactionId}) => [
      {type: TracetestApiTags.TRANSACTION_RUN, id: `${transactionId}-LIST`},
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
  }),

  deleteTransactionRunById: builder.mutation<TTransactionRun, {transactionId: string; runId: string}>({
    query: ({transactionId, runId}) => ({
      url: `/transactions/${transactionId}/run/${runId}`,
      method: HTTP_METHOD.DELETE,
    }),
    invalidatesTags: (result, error, {transactionId}) => [
      {type: TracetestApiTags.TRANSACTION_RUN, id: `${transactionId}-LIST`},
    ],
  }),
});

export default TransactionRunEndpoint;
