import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import Transaction, {TRawTransaction, TRawTransactionResource} from 'models/Transaction.model';
import TransactionService from 'services/Transaction.service';
import {TDraftTransaction} from 'types/Transaction.types';
import TraceTestAPI from '../Tracetest.api';

const defaultHeaders = {'content-type': 'application/json', 'X-Tracetest-Augmented': 'true'};

const transactionEndpoints = TraceTestAPI.injectEndpoints({
  endpoints: builder => ({
    createTransaction: builder.mutation<Transaction, TDraftTransaction>({
      query: transaction => ({
        url: '/transactions',
        method: HTTP_METHOD.POST,
        body: TransactionService.getRawFromDraft(transaction),
        headers: defaultHeaders,
      }),
      transformResponse: (rawTransaction: TRawTransactionResource) => Transaction(rawTransaction),
      invalidatesTags: [
        {type: TracetestApiTags.TRANSACTION, id: 'LIST'},
        {type: TracetestApiTags.RESOURCE, id: 'LIST'},
      ],
    }),
    editTransaction: builder.mutation<Transaction, {transaction: TDraftTransaction; transactionId: string}>({
      query: ({transactionId, transaction}) => ({
        url: `/transactions/${transactionId}`,
        method: HTTP_METHOD.PUT,
        body: TransactionService.getRawFromDraft(transaction),
        headers: defaultHeaders,
      }),
      invalidatesTags: [{type: TracetestApiTags.TRANSACTION, id: 'LIST'}],
    }),
    deleteTransactionById: builder.mutation<Transaction, {transactionId: string}>({
      query: ({transactionId}) => ({
        url: `/transactions/${transactionId}`,
        method: HTTP_METHOD.DELETE,
        headers: defaultHeaders,
      }),
      invalidatesTags: [
        {type: TracetestApiTags.TRANSACTION, id: 'LIST'},
        {type: TracetestApiTags.RESOURCE, id: 'LIST'},
      ],
    }),
    getTransactionById: builder.query<Transaction, {transactionId: string}>({
      query: ({transactionId}) => ({
        url: `/transactions/${transactionId}`,
        headers: defaultHeaders,
      }),
      providesTags: result => [{type: TracetestApiTags.TRANSACTION, id: result?.id}],
      transformResponse: (rawTransaction: TRawTransactionResource) => Transaction(rawTransaction),
    }),
    getTransactionVersionById: builder.query<Transaction, {transactionId: string; version: number}>({
      query: ({transactionId, version}) => `/transactions/${transactionId}/version/${version}`,
      providesTags: result => [{type: TracetestApiTags.TRANSACTION, id: `${result?.id}-${result?.version}`}],
      transformResponse: (rawTest: TRawTransaction) => Transaction.FromRawTransaction(rawTest),
      keepUnusedDataFor: 10,
    }),
  }),
});

export const {
  useCreateTransactionMutation,
  useGetTransactionByIdQuery,
  useDeleteTransactionByIdMutation,
  useEditTransactionMutation,
  useGetTransactionVersionByIdQuery,
  useLazyGetTransactionVersionByIdQuery,
} = transactionEndpoints;
