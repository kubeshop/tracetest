import {difference, sortBy} from 'lodash';
import {useMemo} from 'react';
import {SpanAssertionResult} from '../../types/Assertion.types';

export interface TParsedAssertion {
  key: string;
  spanLabels: string[];
  property: string;
  comparison: string;
  value: string;
  actualValue: string;
  hasPassed: boolean;
  spanId: string;
}

export function useParserAssertionListMemo(
  spanListAssertionResult: SpanAssertionResult[],
  selectorValueList: string[]
): TParsedAssertion[] {
  return useMemo(() => {
    const spanAssertionList = spanListAssertionResult.reduce<Array<TParsedAssertion>>((list, {resultList, span}) => {
      const subResultList = resultList.map<TParsedAssertion>(
        ({propertyName, comparisonValue, operator, actualValue, hasPassed, spanId}) => {
          const spanLabelList = span.signature.map(({value}) => value).concat([`#${spanId.slice(-4)}`]) || [];

          return {
            spanLabels: difference(spanLabelList, selectorValueList),
            spanId,
            key: `${propertyName}-${spanId}`,
            property: propertyName,
            comparison: operator,
            value: comparisonValue,
            actualValue,
            hasPassed,
          };
        }
      );

      return list.concat(subResultList);
    }, []);

    return sortBy(spanAssertionList, ({spanLabels}) => spanLabels.join(''));
  }, [selectorValueList, spanListAssertionResult]);
}
