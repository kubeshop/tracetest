import {TResolveExpressionContext} from 'types/Expression.types';
import {useCallback} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import Validator from 'utils/Validator';

const {useParseExpressionMutation} = TracetestAPI.instance;

const placeholders = {
  var: '${var:',
  env: '${env:',
};

const useValidateUrl = (context: TResolveExpressionContext = {}) => {
  const [parseExpressionMutation, {isLoading}] = useParseExpressionMutation();

  const onParse = useCallback(
    async (url: string) => {
      try {
        const [parsedExpression = ''] = await parseExpressionMutation({
          expression: url,
          context,
        }).unwrap();

        return parsedExpression;
      } catch (e) {
        return Promise.reject(e);
      }
    },
    [context, parseExpressionMutation]
  );

  const onValidate = useCallback(
    async (raw: string): Promise<boolean> => {
      if (!raw) return Promise.resolve(false);

      try {
        const parsed = await onParse(raw);

        const isValid = Validator.url(parsed);

        if (raw.includes(placeholders.var) || raw.includes(placeholders.env)) {
          return Promise.resolve(true);
        }

        return Promise.resolve(isValid);
      } catch (e) {
        return Promise.resolve(false);
      }
    },
    [onParse]
  );

  return {onValidate, isLoading};
};

export default useValidateUrl;
