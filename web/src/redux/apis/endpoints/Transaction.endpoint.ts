import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import Transaction from 'models/Transaction.model';
import Environment from 'models/Environment.model';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {TDraftTransaction, TRawTransaction, TTransaction, TTransactionRun} from 'types/Transaction.types';

const TransactionEndpoint = (builder: TTestApiEndpointBuilder) => ({
  createTransaction: builder.mutation<TTransaction, TDraftTransaction>({
    query: transaction => ({
      url: '/transactions',
      method: HTTP_METHOD.POST,
      body: transaction,
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
      body: transaction,
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
  getTransactionRunById: builder.query<TTransactionRun, {transactionId: string; runId?: string}>({
    query: () => `/tests`,
    providesTags: result => [{type: TracetestApiTags.TRANSACTION, id: result?.id}],
    transformResponse: () => ({
      id: '1',
      transactionVersion: 1,
      results: [],
      environment: Environment({
        name: 'mock',
        id: '1',
        description: 'mock',
        values: [
          {
            key: 'HOST',
            value: 'http://localhost',
          },
          {
            key: 'PORT',
            value: '3000',
          },
        ],
      }),
    }),
  }),
  getTransactionById: builder.query<TTransaction, {transactionId: string}>({
    query: ({transactionId}) => `/transactions/${transactionId}`,
    providesTags: result => [{type: TracetestApiTags.TRANSACTION, id: result?.id}],
    transformResponse: (rawTransaction: TRawTransaction) => Transaction(rawTransaction),
  }),
});

export default TransactionEndpoint;
