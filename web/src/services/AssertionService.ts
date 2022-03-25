import jemspath from 'jmespath';

import {Assertion, COMPARE_OPERATOR, ITrace, LOCATION_NAME} from 'types';

export const assertion1: Assertion = {
  assertionId: 'ABC',
  selectors: [
    {
      locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
      propertyName: 'rpc.system',
      value: 'grpc',
      valueType: 'stringValue',
    },
  ],
  spanAssertions: [
    {
      spanAssertionId: 'ABC',
      locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
      propertyName: 'rpc.grpc.status_code',
      valueType: 'intValue',
      comparisonValue: '0',
      operator: COMPARE_OPERATOR.EQUALS,
    },
  ],
};

export const assertion2: Assertion = {
  assertionId: 'ABC',
  selectors: [
    {
      locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
      propertyName: 'rpc.system',
      value: 'grpc',
      valueType: 'stringValue',
    },
  ],
  spanAssertions: [
    {
      spanAssertionId: 'ABC',
      locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
      propertyName: 'rpc.grpc.status_code',
      valueType: 'intValue',
      comparisonValue: '1',
      operator: COMPARE_OPERATOR.EQUALS,
    },
  ],
};

export const assertion3: Assertion = {
  assertionId: 'ABC',
  selectors: [
    {
      locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
      propertyName: 'db.system',
      value: 'redis',
      valueType: 'stringValue',
    },
  ],
  spanAssertions: [
    {
      spanAssertionId: 'ABC',
      locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
      propertyName: 'db.statement',
      valueType: 'stringValue',
      comparisonValue: 'get key',
      operator: COMPARE_OPERATOR.EQUALS,
    },
  ],
};

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

const buildSelector = (locationName: LOCATION_NAME, query: string) => {
  switch (locationName) {
    case LOCATION_NAME.RESOURCE_ATTRIBUTES:
      return `resource[?${query}]`;
    case LOCATION_NAME.INSTRUMENTATION_LIBRARY:
      return `resource[?${query}]`;
    case LOCATION_NAME.SPAN_ATTRIBUTES:
      return `instrumentationLibrarySpans[?spans[?attributes[?${query}]]]`;
    case LOCATION_NAME.SPAN:
      return `instrumentationLibrarySpans[?spans[?attributes[?${query}]]]`;
    default:
      return '';
  }
};

const flattenTraceArray = () => ` [][instrumentationLibrarySpans[].spans]|[][][]`;

export const runTestAssertion = async (trace: ITrace, assertion: Assertion) => {
  const selectorsMap = assertion.selectors.reduce<{[key in LOCATION_NAME]?: string[]}>((acc, item) => {
    const selectorArray = acc[item.locationName] || [];
    const keySelector = ` key == \`${item.propertyName}\``;
    const valueSelector = buildValueSelector(item.value, '==', item.valueType);
    selectorArray.push(keySelector);
    selectorArray.push(valueSelector);
    acc[item.locationName] = selectorArray;
    return acc;
  }, {});

  const combineSelectors = Object.keys(selectorsMap)
    .map(key => {
      return buildSelector(key as LOCATION_NAME, ` ${selectorsMap[key as keyof typeof LOCATION_NAME]!.join(' && ')}`);
    })
    .join(' && ');

  const itemSelector = `[? ${combineSelectors}]`;

  const selectedSpans = jemspath.search(trace, `resourceSpans|[]| ${itemSelector}`);

  const flattenSelectedSpan = jemspath.search(selectedSpans, flattenTraceArray());

  const assertionTestResultArray = assertion.spanAssertions.map(el => {
    const valueSelector = buildValueSelector(el.comparisonValue, getOperator(el.operator), el.valueType);

    const selector = buildSelector(el.locationName, `key == \`${el.propertyName}\` && ${valueSelector}`);

    const passedSpans = jemspath.search(selectedSpans, `[?${selector}] | ${flattenTraceArray()}`);

    const hasPassed = passedSpans.length === selectedSpans.length;
    const passSpansIds = passedSpans.map((span: any) => span.spanId);

    const failedSpanArray = flattenSelectedSpan.filter(
      (selectedSpan: any) => passSpansIds.indexOf(selectedSpan.spanId) === -1
    );

    return {
      ...el,
      hasPassed,
      spanCount: selectedSpans.length,
      passedSpanCount: passedSpans.length,
      failedSpans: failedSpanArray,
    };
  });

  return assertionTestResultArray;
};
