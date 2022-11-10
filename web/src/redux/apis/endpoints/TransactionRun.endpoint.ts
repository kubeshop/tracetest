import {HTTP_METHOD} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {PaginationResponse} from 'hooks/usePagination';
import TransactionRun from 'models/TransactionRun.model';
import {/* TRawTest, */ TTestApiEndpointBuilder} from 'types/Test.types';
import {TRawTransactionRun, TTransactionRun} from 'types/TransactionRun.types';
// import {TRawTestRun} from '../../../types/TestRun.types';
// import {TRawEnvironment} from '../../../types/Environment.types';

function getTotalCountFromHeaders(meta: any) {
  return Number(meta?.response?.headers.get('x-total-count') || 0);
}

const TransactionRunEndpoint = (builder: TTestApiEndpointBuilder) => ({
  getTransactionRuns: builder.query<
    PaginationResponse<TTransactionRun>,
    {transactionId: string; take?: number; skip?: number}
  >({
    // query: ({transactionId, take = 25, skip = 0}) => `/transactions/${transactionId}/run?take=${take}&skip=${skip}`,
    query: ({transactionId, take = 25, skip = 0}) => `/tests`,
    providesTags: (result, error, {transactionId}) => [
      {type: TracetestApiTags.TRANSACTION_RUN, id: `${transactionId}-LIST`},
    ],
    /* transformResponse: (rawTransactionRuns: TRawTransactionRun[], meta) => ({
      total: getTotalCountFromHeaders(meta),
      items: rawTransactionRuns.map(rawTransactionRun => TransactionRun(rawTransactionRun)),
    }), */
    transformResponse: (rawTransactionRuns: TRawTransactionRun[], meta) => {
      const rawItem = {
        id: '1',
        createdAt: '2022-11-09T17:38:46.165444Z',
        completedAt: '2022-11-09T17:38:46.165444Z',
        state: 'CREATED' as const,
        steps: [],
        stepRuns: [],
        // environment?: TRawEnvironment;
        // metadata?: {[key: string]: string};
      };
      return {
        total: 1,
        items: [TransactionRun(rawItem)],
      };
    },
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
