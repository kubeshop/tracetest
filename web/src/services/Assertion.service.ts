import {search} from 'jmespath';
import {SpanAttributeType} from '../constants/SpanAttribute.constants';
import {IAssertion, IAssertionResult, IItemSelector, ISpanAssertionResult} from '../types/Assertion.types';
import {ISpan} from '../types/Span.types';
import {ITrace} from '../types/Trace.types';
import {escapeString, isJson} from '../utils/Common';
import OperatorService from './Operator.service';

const buildValueSelector = (comparisonValue: string, compareOperator: string, type: string) => {
  if (compareOperator === 'contains') return `contains(value, \`${comparisonValue}\`)`;

  if ([SpanAttributeType.intValue, SpanAttributeType.doubleValue].includes(type as SpanAttributeType)) {
    return `to_number(value) ${compareOperator} \`${escapeString(comparisonValue)}\``;
  }

  if (isJson(comparisonValue)) {
    return `value ${compareOperator} \`${JSON.stringify(comparisonValue)}\``;
  }

  return `value ${compareOperator} \`${escapeString(comparisonValue)}\``;
};

const buildSelector = (conditions: string[]) => `${conditions.map(cond => `${cond}`).join(' && ')}`;

const getSelectorList = (itemSelectors: IItemSelector[] = []) => {
  return itemSelectors.map<string>(({propertyName, value, valueType}) => {
    const keySelector = ` key == \`${propertyName}\``;
    const valueSelector = buildValueSelector(value, '==', valueType);
    return `${[keySelector, valueSelector]!.join(' && ')}`;
  }, {});
};

const AssertionService = () => ({
  runBySpan(span: ISpan, {spanAssertions = [], selectors}: IAssertion): Array<ISpanAssertionResult> {
    const {spanId, attributes} = span;
    const itemSelector = getSelectorList(selectors || []);
    const itMatches = !selectors || search(span.attributeList, escapeString(`[? ${itemSelector.join('] && [? ')}]`))[0];

    if (!itMatches) return [];

    return spanAssertions.map(spanAssertion => {
      const {comparisonValue, operator, propertyName, valueType} = spanAssertion;
      const valueSelector = buildValueSelector(comparisonValue, OperatorService.getOperatorSymbol(operator), valueType);

      const selector = `${buildSelector([`[? key == \`${propertyName}\` && ${valueSelector}]`])}`;

      const [hasPassed] = search(span.attributeList, selector);

      return {
        ...spanAssertion,
        spanId,
        hasPassed: Boolean(hasPassed),
        actualValue: attributes[propertyName]?.value || '',
      };
    });
  },

  runByTrace(trace: ITrace, assertion: IAssertion): IAssertionResult {
    const itemSelector = `${getSelectorList(assertion.selectors || [])
      .map(condition => `attributeList[? ${condition}]`)
      .join(' && ')}`;

    const spanList: ISpan[] = assertion.selectors
      ? search(trace.spans, escapeString(`[? ${itemSelector}]`)) || []
      : trace.spans;

    return {
      assertion,
      spanListAssertionResult: spanList.map(span => ({
        span,
        resultList: this.runBySpan(span, assertion),
      })),
    };
  },

  getEffectedSpansCount(trace: ITrace, selectors: IItemSelector[]) {
    if (selectors.length === 0) return trace.spans;

    const itemSelector = `${getSelectorList(selectors)
      .map(condition => `attributeList[? ${condition}]`)
      .join(' && ')}`;
    const spanList: ISpan[] = search(trace.spans, escapeString(`[? ${itemSelector}]`)) || [];

    return spanList;
  },

  newSelectorLogic(selectorList: IItemSelector[]): string {
    function getValue({value, valueType}: {value: string; valueType: string}): string {
      let result = ``;
      // add quotes if value is a string
      if (valueType === 'stringValue') result += `'`;
      result += value;
      // add quotes if value is a string
      if (valueType === 'stringValue') result += `'`;
      return result;
    }

    const getFilters = (selectors: IItemSelector[]) => {
      return selectors.map(selector => {
        const {propertyName, operator, value, valueType} = selector;
        return `${propertyName}${operator ? ` ${operator.toLowerCase()} ` : '='}${getValue({value, valueType})}`;
      });
    };
    return selectorList.length === 0 ? '' : `span[${getFilters(selectorList).join(' ')}]`;
  },
});

export default AssertionService();
