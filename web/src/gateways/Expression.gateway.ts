import {endpoints} from 'redux/apis/TraceTest.api';
import {TParseRequestInfo} from 'types/Expression.types';

const {parseExpression} = endpoints;

const ExpressionGateway = () => ({
  parseExpression(props: TParseRequestInfo) {
    return parseExpression.initiate(props);
  },
});

export default ExpressionGateway();
