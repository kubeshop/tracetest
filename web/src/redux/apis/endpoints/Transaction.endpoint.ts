import TransactionMock from 'models/__mocks__/Transaction.mock';
import {TTransaction} from 'types/Transaction.types';
import {TTestApiEndpointBuilder} from 'types/Test.types';
import {TracetestApiTags} from 'constants/Test.constants';

const TransactionEndpoint = (builder: TTestApiEndpointBuilder) => ({
  deleteTransactionById: builder.mutation<TTransaction, {transactionId: string}>({
    query: ({transactionId}) => ({url: `/transactions/${transactionId}`, method: 'DELETE'}),
    invalidatesTags: [{type: TracetestApiTags.TRANSACTION, id: 'LIST'}],
  }),
  getTransactionRunById: builder.query<TTransaction, {transactionId: string; runId?: string}>({
    query: () => `/tests`,
    providesTags: result => [{type: TracetestApiTags.TRANSACTION, id: result?.id}],
    transformResponse: () => TransactionMock.model(),
  }),
  getTransactionById: builder.query<TTransaction, {transactionId: string}>({
    query: () => `/tests`,
    providesTags: result => [{type: TracetestApiTags.TRANSACTION, id: result?.id}],
    transformResponse: () => TransactionMock.model(),
  }),
});

export default TransactionEndpoint;
