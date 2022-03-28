import jemspath from 'jmespath';

import {Assertion, AssertionResult, COMPARE_OPERATOR, ItemSelector, ITrace, LOCATION_NAME} from 'types';

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

const buildSelector = (locationName: LOCATION_NAME, conditions: string[]) => {
  switch (locationName) {
    case LOCATION_NAME.RESOURCE_ATTRIBUTES:
      return `${conditions.map(cond => `resource.attributes[?${cond}]`).join(' && ')}`;
    case LOCATION_NAME.INSTRUMENTATION_LIBRARY:
      return `${conditions.map(cond => `resource.attributes[?${cond}]`).join(' && ')}`;
    case LOCATION_NAME.SPAN_ATTRIBUTES:
      return `instrumentationLibrarySpans[?spans[?${conditions.map(cond => `attributes[?${cond}]`).join(' && ')}]]`;
    case LOCATION_NAME.SPAN:
      return `instrumentationLibrarySpans[?spans[?${conditions.map(cond => `attributes[?${cond}]`).join(' && ')}]]`;
    default:
      return '';
  }
};

const flattenTraceArray = () => ` [][instrumentationLibrarySpans[].spans]|[][][]`;

export const buildItemSelectorQuery = (itemSelectors: ItemSelector[]) => {
  const selectorsMap = itemSelectors.reduce<string[]>((acc, item) => {
    const keySelector = ` key == \`${item.propertyName}\``;
    const valueSelector = buildValueSelector(item.value, '==', item.valueType);

    acc.push(` ${[keySelector, valueSelector]!.join(' && ')}`);
    return acc;
  }, []);

  const itemSelector = `[? ${buildSelector(LOCATION_NAME.SPAN_ATTRIBUTES, selectorsMap)}]`;
  return itemSelector;
};

export const runTestAssertion = (trace: ITrace, assertion: Assertion): AssertionResult[] => {
  if (!assertion?.selectors) {
    return [];
  }

  const itemSelector = buildItemSelectorQuery(assertion.selectors);
  const selectedSpans = jemspath.search(trace, `resourceSpans|[]| ${itemSelector}`);
  const flattenSelectedSpan = jemspath.search(selectedSpans, flattenTraceArray());

  const assertionTestResultArray = assertion.spanAssertions.map(el => {
    const valueSelector = buildValueSelector(el.comparisonValue, getOperator(el.operator), el.valueType);

    const selector = buildSelector(el.locationName, [`key == \`${el.propertyName}\` && ${valueSelector}`]);
    console.log('@@selector', selector);
    const passedSpans = jemspath.search(selectedSpans, `[?${selector}] | ${flattenTraceArray()}`);
    console.log('@@passedSpans', passedSpans);
    const hasPassed = passedSpans.length === selectedSpans.length;
    const passSpansIds = passedSpans.map((span: any) => span.spanId);

    const failedSpanArray = flattenSelectedSpan.filter(
      (selectedSpan: any) => passSpansIds.indexOf(selectedSpan.spanId) === -1
    );

    return {
      ...el,
      selector: assertion.selectors.map(item => item.propertyName).join(' '),
      hasPassed,
      spanCount: selectedSpans.length,
      passedSpanCount: passedSpans.length || 0,
      failedSpans: failedSpanArray,
    };
  });

  return assertionTestResultArray;
};
