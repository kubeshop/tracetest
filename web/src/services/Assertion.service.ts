import {search} from 'jmespath';
import {escapeString} from '../utils/Common';
import OperatorService from './Operator.service';
import {ISpan} from '../types/Span.types';
import {ITrace} from '../types/Trace.types';
import {IAssertion, IAssertionResult, IItemSelector, ISpanAssertionResult} from '../types/Assertion.types';
import {SpanAttributeType} from '../constants/SpanAttribute.constants';

const buildValueSelector = (comparisonValue: string, compareOperator: string, type: string) => {
  if ([SpanAttributeType.intValue, SpanAttributeType.doubleValue].includes(type as SpanAttributeType)) {
    return `to_number(value) ${compareOperator} \`${comparisonValue}\``;
  }

  if (compareOperator === 'contains') return `contains(value, \`${comparisonValue}\`)`;

  return `value ${compareOperator} \`${comparisonValue}\``;
};

const buildSelector = (conditions: string[]) => `${conditions.map(cond => `${cond}`).join(' && ')}`;

const getSelectorList = (itemSelectors: IItemSelector[]) => {
  const selectorList = itemSelectors.map<string>(({propertyName, value, valueType}) => {
    const keySelector = ` key == \`${propertyName}\``;
    const valueSelector = buildValueSelector(value, '==', valueType);
    const condition = `${[keySelector, valueSelector]!.join(' && ')}`;

    return condition;
  }, {});

  return selectorList;
};

const AssertionService = () => ({
  runBySpan(span: ISpan, {spanAssertions = [], selectors}: IAssertion): Array<ISpanAssertionResult> {
    const {spanId, attributes} = span;
    const itemSelector = getSelectorList(selectors);
    const [itMatches] = search(span.attributeList, escapeString(`[? ${itemSelector.join('] && [? ')}]`));

    if (!itMatches) return [];

    const assertionTestResultArray = spanAssertions.map(spanAssertion => {
      const {comparisonValue, operator, propertyName, valueType} = spanAssertion;
      const valueSelector = buildValueSelector(comparisonValue, OperatorService.getOperatorSymbol(operator), valueType);

      const selector = `${buildSelector([`[? key == \`${propertyName}\` && ${valueSelector}]`])}`;

      const [hasPassed] = search(span.attributeList, escapeString(selector));

      return {
        ...spanAssertion,
        spanId,
        hasPassed: Boolean(hasPassed),
        actualValue: attributes[propertyName]?.value || '',
      };
    });

    return assertionTestResultArray;
  },

  runByTrace(trace: ITrace, assertion: IAssertion): IAssertionResult {
    if (!assertion?.selectors) return {assertion, spanListAssertionResult: []};

    const itemSelector = `${getSelectorList(assertion.selectors)
      .map(condition => `attributeList[? ${condition}]`)
      .join(' && ')}`;
    const spanList: ISpan[] = search(trace.spans, escapeString(`[? ${itemSelector}]`)) || [];

    return {
      assertion,
      spanListAssertionResult: spanList.map(span => ({
        span,
        resultList: this.runBySpan(span, assertion),
      })),
    };
  },

  getEffectedSpansCount(trace: ITrace, selectors: IItemSelector[]) {
    if (selectors.length === 0) return 0;

    const itemSelector = `${getSelectorList(selectors)
      .map(condition => `attributeList[? ${condition}]`)
      .join(' && ')}`;
    const spanList: ISpan[] = search(trace.spans, escapeString(`[? ${itemSelector}]`)) || [];

    return spanList.length;
  },
});

export default AssertionService();
