import {useCallback, useEffect, useState} from 'react';
import {TResolveExpressionContext, TResolveRequestInfo} from 'types/Expression.types';
import {useParseExpressionMutation} from 'redux/apis/TraceTest.api';

const useTooltip = (context: TResolveExpressionContext = {}) => {
  const [parseExpressionMutation] = useParseExpressionMutation();
  const [prevExpression, setPrevExpression] = useState<string>('');
  const [prevRawExpression, setPrevRawExpression] = useState<string>('');

  useEffect(() => {
    if (prevExpression) parseExpression({expression: prevExpression, context});
  }, [context.environmentId, prevExpression]);

  const parseExpression = useCallback(
    async (props: TResolveRequestInfo) => {
      const parsedExpression = await parseExpressionMutation(props).unwrap();

      setPrevExpression(parsedExpression);
      setPrevRawExpression(props.expression || '');
    },
    [parseExpressionMutation]
  );

  const onHover = useCallback(
    (rawExpression: string) => {
      const isSameAsPrev = prevRawExpression === rawExpression;

      if (isSameAsPrev) return;

      parseExpression({
        expression: rawExpression,
        context,
      });
    },
    [context, parseExpression, prevRawExpression]
  );

  return {onHover, expression: prevExpression};
};

export default useTooltip;
