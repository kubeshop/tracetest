import {HTTP_METHOD} from 'constants/Common.constants';
import AssertionService from 'services/Assertion.service';
import {TResolveRequestInfo, TResolveResponseInfo} from 'types/Expression.types';
import { TTestApiEndpointBuilder } from '../Tracetest.api';

export const expressionEndpoints = (builder: TTestApiEndpointBuilder) => ({
  parseExpression: builder.mutation<string[], TResolveRequestInfo>({
    query: ({expression, context: {spanId = '', runId = 0, variableSetId = '', testId = '', selector = ''} = {}}) => ({
      url: '/expressions/resolve',
      method: HTTP_METHOD.POST,
      body: {
        expression: AssertionService.extractExpectedString(expression),
        context: {spanId, runId, testId, selector, variableSetId},
      },
    }),
    transformResponse: ({resolvedValues = []}: TResolveResponseInfo) => resolvedValues,
  }),
});
