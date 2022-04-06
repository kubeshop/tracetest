import {search} from 'jmespath';

import {
  Assertion,
  AssertionResult,
  COMPARE_OPERATOR,
  ItemSelector,
  ITrace,
  LOCATION_NAME,
  ResourceSpan,
  SpanAssertionResult,
} from 'types';
import { getSpanValue } from './SpanService';

const getOperator = (op: COMPARE_OPERATOR) => {
  switch (op) {
    case COMPARE_OPERATOR.EQUALS:
      return '==';
    case COMPARE_OPERATOR.NOTEQUALS:
      return '!=';
    case COMPARE_OPERATOR.GREATERTHAN:
      return '>';
    case COMPARE_OPERATOR.LESSTHAN:
      return '<';
    default:
      return '==';
  }
};

const buildValueSelector = (comparisonValue: string, compareOperator: string, valueType: string) => {
  if (valueType === 'intValue') {
    return `to_number(value.${valueType}) ${compareOperator} \`${comparisonValue}\``;
  }
  return `value.${valueType} ${compareOperator} \`${comparisonValue}\``;
};

const buildSelector = (locationName: LOCATION_NAME, conditions: string[], spanId?: string) => {
  switch (locationName) {
    case LOCATION_NAME.INSTRUMENTATION_LIBRARY:
    case LOCATION_NAME.RESOURCE_ATTRIBUTES:
      return `${conditions.map(cond => `resource.attributes[?${cond}]`).join(' && ')}`;
    case LOCATION_NAME.SPAN_ATTRIBUTES:
    case LOCATION_NAME.SPAN:
      return `instrumentationLibrarySpans[?spans[?${conditions.map(cond => `attributes[?${cond}]`).join(' && ')}]]`;
    case LOCATION_NAME.SPAN_ID:
      return `instrumentationLibrarySpans[?spans[?spanId == \`${spanId}\` && ${conditions
        .map(cond => `attributes[?${cond}]`)
        .join(' && ')}]]`;
    default:
      return '';
  }
};

const buildConditionArray = (itemSelectors: ItemSelector[]) => {
  const selectorsMap = itemSelectors.reduce<string[]>((acc, item) => {
    const keySelector = ` key == \`${item.propertyName}\``;
    const valueSelector = buildValueSelector(item.value, '==', item.valueType);

    acc.push(` ${[keySelector, valueSelector]!.join(' && ')}`);
    return acc;
  }, []);

  return selectorsMap;
};

export const buildItemSelectorQuery = (itemSelectors: ItemSelector[]) => {
  const selectorsMap = buildConditionArray(itemSelectors);

  const itemSelector = `[? ${buildSelector(LOCATION_NAME.SPAN_ATTRIBUTES, selectorsMap)}]`;
  return itemSelector;
};

export const runSpanAssertion = (span: ResourceSpan, assertion: Assertion): Array<SpanAssertionResult> => {
  const assertionTestResultArray = assertion.spanAssertions.map(spanAssertion => {
    const {comparisonValue, operator, valueType, locationName, propertyName} = spanAssertion;
    const valueSelector = buildValueSelector(comparisonValue, getOperator(operator), valueType);

    const selector = `${buildSelector(locationName, [`key == \`${propertyName}\` && ${valueSelector}`])}`;
    const [passedSpan] = search(span, selector);

    return {
      ...spanAssertion,
      hasPassed: Boolean(passedSpan),
      actualValue: getSpanValue(span, locationName, valueType, propertyName),
    };
  });

  return assertionTestResultArray;
};

export const runAssertionBySpanId = (
  spanId: string,
  trace: ITrace,
  assertion: Assertion
): SpanAssertionResult[] | undefined => {
  if (!assertion.selectors) return undefined;
  const conditionList = buildConditionArray(assertion.selectors);
  const itemSelector = `[? ${buildSelector(LOCATION_NAME.SPAN_ID, conditionList, spanId)}]`;
  const [span]: Array<ResourceSpan> = search(trace, `resourceSpans|[]| ${itemSelector}`);

  if (!span) return undefined;

  return runSpanAssertionByResourceSpan(span, assertion);
};

export const runSpanAssertionByResourceSpan = (span: ResourceSpan, assertion: Assertion): SpanAssertionResult[] => {
  const spanResult = runSpanAssertion(span, assertion);

  return spanResult;
};

export const runAssertionByTrace = (trace: ITrace, assertion: Assertion): AssertionResult => {
  if (!assertion?.selectors) return {assertion, spanListAssertionResult: []};

  const itemSelector = buildItemSelectorQuery(assertion.selectors);
  const spanList: Array<ResourceSpan> = search(trace, `resourceSpans|[]| ${itemSelector}`);

  return {
    assertion,
    spanListAssertionResult: spanList.map(span => ({
      span,
      resultList: runSpanAssertionByResourceSpan(span, assertion),
    })),
  };
};
