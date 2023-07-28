import {HTTP_METHOD} from 'constants/Common.constants';
import AssertionService from 'services/Assertion.service';
import {TResolveRequestInfo, TResolveResponseInfo} from 'types/Expression.types';
import TraceTestAPI from '../Tracetest.api';

const expressionEndpoints = TraceTestAPI.injectEndpoints({
  endpoints: builder => ({
    parseExpression: builder.mutation<string[], TResolveRequestInfo>({
      query: ({
        expression,
        context: {spanId = '', runId = '', environmentId = '', testId = '', selector = ''} = {},
      }) => ({
        url: '/expressions/resolve',
        method: HTTP_METHOD.POST,
        body: {
          expression: AssertionService.extractExpectedString(expression),
          context: {spanId, runId, testId, selector, environmentId},
        },
      }),
      transformResponse: ({resolvedValues = []}: TResolveResponseInfo) => resolvedValues,
    }),
  }),
});

export const {useParseExpressionMutation} = expressionEndpoints;
