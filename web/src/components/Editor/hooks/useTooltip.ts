import {useCallback, useState} from 'react';
import {TResolveExpressionContext, TResolveRequestInfo} from 'types/Expression.types';
import {useParseExpressionMutation} from 'redux/apis/TraceTest.api';

const useTooltip = (context: TResolveExpressionContext = {}) => {
  const [parseExpressionMutation] = useParseExpressionMutation();
  const [prevExpression, setPrevExpression] = useState<string[]>([]);
  const [prevRawExpression, setPrevRawExpression] = useState<string>('');
  const [prevContext, setPrevContext] = useState<TResolveExpressionContext>({});

  const parseExpression = useCallback(
    async (props: TResolveRequestInfo) => {
      const isSameAsPrev =
        prevRawExpression === props.expression && JSON.stringify(prevContext) === JSON.stringify(context);

      if (isSameAsPrev) return prevExpression;

      const parsedExpression = await parseExpressionMutation(props).unwrap();

      setPrevExpression(parsedExpression);
      setPrevContext(props.context || {});
      setPrevRawExpression(props.expression || '');
    },
    [context, parseExpressionMutation, prevContext, prevExpression, prevRawExpression]
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

  return {onHover, resolvedValues: prevExpression};
};

export default useTooltip;
