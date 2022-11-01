import {TTestApiEndpointBuilder} from 'types/Test.types';
import {TParseRequestInfo, TParseResponseInfo} from 'types/Expression.types';
import {HTTP_METHOD} from 'constants/Common.constants';
import AssertionService from 'services/Assertion.service';

const ExpressionEndpoint = (builder: TTestApiEndpointBuilder) => ({
  parseExpression: builder.mutation<string, TParseRequestInfo>({
    query: ({expression, context: {spanId = '', runId = '', environmentId = '', testId = '', selector = ''} = {}}) => ({
      url: '/expressions/parse',
      method: HTTP_METHOD.POST,
      body: {
        expression: AssertionService.extractExpectedString(expression),
        context: {spanId, runId, testId, selector, environmentId},
      },
    }),
    transformResponse: (res: TParseResponseInfo) => res.expression || '',
  }),
});

export default ExpressionEndpoint;
