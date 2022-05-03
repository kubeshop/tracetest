import {search} from 'jmespath';
import {escapeString} from '../utils/Common';
import OperatorService from './Operator.service';
import {TSpan} from '../types/Span.types';
import {TTrace} from '../types/Trace.types';
import {TAssertion, TAssertionResult, TItemSelector, TSpanAssertionResult} from '../types/Assertion.types';

const buildValueSelector = (comparisonValue: string, compareOperator: string, type: string) => {
  if (compareOperator === 'contains') return `contains(value, \`${comparisonValue}\`)`;

  if (['intValue', 'doubleValue'].includes(type)) {
    return `to_number(value) ${compareOperator} \`${comparisonValue}\``;
  }

  return `value ${compareOperator} \`${comparisonValue}\``;
};

const buildSelector = (conditions: string[]) => `${conditions.map(cond => `${cond}`).join(' && ')}`;

const getSelectorList = (itemSelectors: TItemSelector[]) => {
  const selectorList = itemSelectors.map<string>(({propertyName, value, valueType}) => {
    const keySelector = ` key == \`${propertyName}\``;
    const valueSelector = buildValueSelector(value!, '==', valueType);
    const condition = `${[keySelector, valueSelector]!.join(' && ')}`;

    return condition;
  }, {});

  return selectorList;
};

const AssertionService = () => ({
  runBySpan(span: TSpan, {spanAssertions = [], selectors}: TAssertion): Array<TSpanAssertionResult> {
    const {spanId, attributes} = span;
    const itemSelector = getSelectorList(selectors);
    const [itMatches] = search(span.attributeList, escapeString(`[? ${itemSelector.join('] && [? ')}]`));

    if (!itMatches) return [];

    const assertionTestResultArray = spanAssertions.map(spanAssertion => {
      const {comparisonValue, operator, propertyName, valueType} = spanAssertion;
      const valueSelector = buildValueSelector(comparisonValue!, OperatorService.getOperatorSymbol(operator), valueType);

      const selector = `${buildSelector([`[? key == \`${propertyName}\` && ${valueSelector}]`])}`;

      const [hasPassed] = search(span.attributeList, escapeString(selector));

      return {
        ...spanAssertion,
        spanId,
        hasPassed: Boolean(hasPassed),
        actualValue: attributes[propertyName!]?.value || '',
      };
    });

    return assertionTestResultArray;
  },

  runByTrace(trace: TTrace, assertion: TAssertion): TAssertionResult {
    if (!assertion?.selectors) return {assertion, spanListAssertionResult: []};

    const itemSelector = `${getSelectorList(assertion.selectors)
      .map(condition => `attributeList[? ${condition}]`)
      .join(' && ')}`;
    const spanList: TSpan[] = search(trace.spans, escapeString(`[? ${itemSelector}]`)) || [];

    return {
      assertion,
      spanListAssertionResult: spanList.map(span => ({
        span,
        resultList: this.runBySpan(span, assertion),
      })),
    };
  },

  getEffectedSpansCount(trace: TTrace, selectors: TItemSelector[]) {
    if (selectors.length === 0) return 0;

    const itemSelector = `${getSelectorList(selectors)
      .map(condition => `attributeList[? ${condition}]`)
      .join(' && ')}`;
    const spanList: TSpan[] = search(trace.spans, escapeString(`[? ${itemSelector}]`)) || [];

    return spanList.length;
  },
});

export default AssertionService();
