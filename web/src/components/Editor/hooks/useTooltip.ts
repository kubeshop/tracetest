import {useCallback, useState} from 'react';
import {useAppDispatch} from 'redux/hooks';
import ExpressionGateway from 'gateways/Expression.gateway';
import {TParseExpressionContext, TParseRequestInfo} from 'types/Expression.types';

const useTooltip = (context: TParseExpressionContext = {}) => {
  const dispatch = useAppDispatch();
  const [prevExpression, setPrevExpression] = useState<string>('');
  const [prevRawExpression, setPrevRawExpression] = useState<string>('');

  const parseExpression = useCallback(
    async (props: TParseRequestInfo) => {
      const isSameAsPrev = prevRawExpression === props.expression;

      if (isSameAsPrev) return prevExpression;

      const parsedExpression = await dispatch(ExpressionGateway.parseExpression(props)).unwrap();

      setPrevExpression(parsedExpression);
      setPrevRawExpression(props.expression || '');
    },
    [dispatch, prevExpression, prevRawExpression]
  );

  const onHover = useCallback(
    (rawExpression: string) => {
      parseExpression({
        expression: rawExpression,
        context,
      });
    },
    [context, parseExpression]
  );

  return {onHover, expression: prevExpression};
};

export default useTooltip;
