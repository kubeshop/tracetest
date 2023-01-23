import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import Transaction from 'models/Transaction.model';
import TransactionService from 'services/Transaction.service';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {TDraftTransaction, TRawTransaction, TTransaction} from 'types/Transaction.types';

const TransactionEndpoint = (builder: TTestApiEndpointBuilder) => ({
  createTransaction: builder.mutation<TTransaction, TDraftTransaction>({
    query: transaction => ({
      url: '/transactions',
      method: HTTP_METHOD.POST,
      body: TransactionService.getRawFromDraft(transaction),
    }),
    transformResponse: (rawTransaction: TRawTransaction) => Transaction(rawTransaction),
    invalidatesTags: [
      {type: TracetestApiTags.TRANSACTION, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  editTransaction: builder.mutation<TTransaction, {transaction: TDraftTransaction; transactionId: string}>({
    query: ({transactionId, transaction}) => ({
      url: `/transactions/${transactionId}`,
      method: HTTP_METHOD.PUT,
      body: TransactionService.getRawFromDraft(transaction),
    }),
    invalidatesTags: [{type: TracetestApiTags.TRANSACTION, id: 'LIST'}],
  }),
  deleteTransactionById: builder.mutation<TTransaction, {transactionId: string}>({
    query: ({transactionId}) => ({url: `/transactions/${transactionId}`, method: HTTP_METHOD.DELETE}),
    invalidatesTags: [
      {type: TracetestApiTags.TRANSACTION, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  getTransactionById: builder.query<TTransaction, {transactionId: string}>({
    query: ({transactionId}) => `/transactions/${transactionId}`,
    providesTags: result => [{type: TracetestApiTags.TRANSACTION, id: result?.id}],
    transformResponse: (rawTransaction: TRawTransaction) => Transaction(rawTransaction),
  }),
  getTransactionVersionById: builder.query<TTransaction, {transactionId: string; version: number}>({
    query: ({transactionId, version}) => `/transactions/${transactionId}/version/${version}`,
    providesTags: result => [{type: TracetestApiTags.TRANSACTION, id: `${result?.id}-${result?.version}`}],
    transformResponse: (rawTest: TRawTransaction) => Transaction(rawTest),
    keepUnusedDataFor: 10,
  }),
});

export default TransactionEndpoint;
