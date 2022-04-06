import {search} from 'jmespath';

import {
  Assertion,
  AssertionResult,
  COMPARE_OPERATOR,
  ItemSelector,
  ITrace,
  LOCATION_NAME,
  ResourceSpan,
  ISpanAttributes,
  SpanAssertionResult,
} from 'types';

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

export const selectSpanValue = (span: ResourceSpan, locationName: LOCATION_NAME, valueType: string, key: string) => {
  switch (locationName) {
    case LOCATION_NAME.INSTRUMENTATION_LIBRARY:
    case LOCATION_NAME.RESOURCE_ATTRIBUTES:
      return search(span, `resource.attributes[? key==\`${key}\`].value.${valueType}`);
    case LOCATION_NAME.SPAN:
    case LOCATION_NAME.SPAN_ATTRIBUTES: {
      const attributeList: ISpanAttributes = search(span, `instrumentationLibrarySpans[].spans[].attributes | []`);

      return attributeList.find(attribute => attribute.key === key)?.value[valueType];
    }
    default:
      return '';
  }
};

const flattenTraceArray = () => ` [][instrumentationLibrarySpans[].spans]|[][][]`;

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

export const runSpanListAssertion = (spanList: Array<ResourceSpan>, assertion: Assertion) => {
  const flattenSelectedSpan = search(spanList, flattenTraceArray());

  const assertionTestResultArray = assertion.spanAssertions.map(el => {
    const valueSelector = buildValueSelector(el.comparisonValue, getOperator(el.operator), el.valueType);

    const selector = buildSelector(el.locationName, [`key == \`${el.propertyName}\` && ${valueSelector}`]);
    const passedSpans = search(spanList, `[?${selector}] | ${flattenTraceArray()}`);
    const hasPassed = passedSpans.length === spanList.length;
    const passSpansIds = passedSpans.map((span: any) => span.spanId);

    const failedSpanArray = flattenSelectedSpan.filter(
      (selectedSpan: any) => passSpansIds.indexOf(selectedSpan.spanId) === -1
    );

    return {
      ...el,
      selector: assertion.selectors.map(item => item.propertyName).join(' '),
      hasPassed,
      spanCount: spanList.length,
      passedSpanCount: passedSpans.length || 0,
      failedSpans: failedSpanArray,
      passSpansIds,
    };
  });

  return assertionTestResultArray;
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
      actualValue: selectSpanValue(span, locationName, valueType, propertyName),
    };
  });

  return assertionTestResultArray;
};

export const runSpanAssertionList = (spanId: string, trace: ITrace, assertion: Assertion) => {
  if (!assertion.selectors) return undefined;
  const conditionList = buildConditionArray(assertion.selectors);
  const itemSelector = `[? ${buildSelector(LOCATION_NAME.SPAN_ID, conditionList, spanId)}]`;
  const [span]: Array<ResourceSpan> = search(trace, `resourceSpans|[]| ${itemSelector}`);

  if (!span) return undefined;

  const spanResult = runSpanAssertion(span, assertion);

  return spanResult;
};

export const runTraceAssertion = (trace: ITrace, assertion: Assertion): AssertionResult[] => {
  if (!assertion?.selectors) {
    return [];
  }

  const itemSelector = buildItemSelectorQuery(assertion.selectors);
  const selectedSpans: Array<ResourceSpan> = search(trace, `resourceSpans|[]| ${itemSelector}`);

  return runSpanListAssertion(selectedSpans, assertion);
};
