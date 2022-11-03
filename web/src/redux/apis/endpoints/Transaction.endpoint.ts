import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import TransactionMock from 'models/__mocks__/Transaction.mock';
import Transaction from 'models/Transaction.model';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {TDraftTransaction, TRawTransaction, TTransaction} from 'types/Transaction.types';

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
  deleteTransactionById: builder.mutation<TTransaction, {transactionId: string}>({
    query: ({transactionId}) => ({url: `/transactions/${transactionId}`, method: 'DELETE'}),
    invalidatesTags: [
      {type: TracetestApiTags.TRANSACTION, id: 'LIST'},
      {type: TracetestApiTags.RESOURCE, id: 'LIST'},
    ],
  }),
  getTransactionRunById: builder.query<TTransaction, {transactionId: string; runId?: string}>({
    query: () => `/tests`,
    providesTags: result => [{type: TracetestApiTags.TRANSACTION, id: result?.id}],
    transformResponse: () => TransactionMock.model(),
  }),
  getTransactionById: builder.query<TTransaction, {transactionId: string}>({
    query: ({transactionId}) => `/transactions/${transactionId}`,
    providesTags: result => [{type: TracetestApiTags.TRANSACTION, id: result?.id}],
    transformResponse: (rawTransaction: TRawTransaction) => Transaction(rawTransaction),
  }),
});

export default TransactionEndpoint;
