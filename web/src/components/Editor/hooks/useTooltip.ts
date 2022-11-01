import {useCallback, useState} from 'react';
import {TParseExpressionContext, TParseRequestInfo} from 'types/Expression.types';
import {useParseExpressionMutation} from 'redux/apis/TraceTest.api';

const useTooltip = (context: TParseExpressionContext = {}) => {
  const [parseExpressionMutation] = useParseExpressionMutation();
  const [prevExpression, setPrevExpression] = useState<string>('');
  const [prevRawExpression, setPrevRawExpression] = useState<string>('');

  const parseExpression = useCallback(
    async (props: TParseRequestInfo) => {
      const isSameAsPrev = prevRawExpression === props.expression;

      if (isSameAsPrev) return prevExpression;

      const parsedExpression = await parseExpressionMutation(props).unwrap();

      setPrevExpression(parsedExpression);
      setPrevRawExpression(props.expression || '');
    },
    [parseExpressionMutation, prevExpression, prevRawExpression]
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
