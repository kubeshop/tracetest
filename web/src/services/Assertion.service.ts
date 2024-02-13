import countBy from 'lodash/countBy';
import uniq from 'lodash/uniq';

import {durationRegExp} from 'constants/Common.constants';
import {CompareOperatorSymbolMap, OperatorRegexp} from 'constants/Operator.constants';
import {Attributes} from 'constants/SpanAttribute.constants';
import {TestSpecErrors} from 'constants/TestSpecs.constants';
import AssertionResult, {TRawAssertionResult} from 'models/AssertionResult.model';
import {ICheckResult, TStructuredAssertion} from 'types/Assertion.types';
import {TCompareOperatorSymbol} from 'types/Operator.types';
import {isJson} from 'utils/Common';

const isNumeric = (num: string): boolean => /^-?\d+(?:\.\d+)?$/.test(num);
const isNumericTime = (num: string): boolean => durationRegExp.test(num);

const AssertionService = () => ({
  extractExpectedString(input?: string): string | undefined {
    if (!input) return input;
    const formatted = input.trim();

    if (isJson(input)) return `'${input}'`;

    if (Object.values(Attributes).includes(formatted)) return formatted;
    if (Object.values(Attributes).some(aa => formatted.includes(aa))) return formatted;
    if (isNumeric(formatted) || isNumericTime(formatted)) return formatted;

    const isQuoted = formatted.startsWith('"') && formatted.endsWith('"');

    return isQuoted ? formatted : this.quotedString(formatted);
  },
  quotedString(str: string): string {
    return `\"${str}\"`;
  },
  getSpanIds(resultList: TRawAssertionResult[]) {
    const spanIds = resultList
      .flatMap(assertion => assertion?.spanResults?.map(span => span.spanId ?? '') ?? [])
      .filter(spanId => Boolean(spanId));
    return uniq(spanIds);
  },

  getTotalPassedChecks(resultList: AssertionResult[]) {
    const passedResults = resultList.flatMap(({spanResults}) => spanResults.map(({passed}) => passed));
    return countBy(passedResults);
  },

  isValidError(error: string) {
    return TestSpecErrors.some(testSpecError => error.includes(testSpecError));
  },

  hasError(resultList: AssertionResult[]) {
    return resultList
      .map(({spanResults}) => spanResults.some(({error}) => !!error && this.isValidError(error)))
      .some(result => !!result);
  },

  getResultsHashedBySpanId(resultList: AssertionResult[], spanIds: string[] = []) {
    return resultList
      .flatMap(({assertion, spanResults}) => spanResults.map(spanResult => ({result: spanResult, assertion})))
      .filter(({result}) => !spanIds.length || spanIds.includes(result.spanId))
      .reduce((prev: Record<string, ICheckResult[]>, curr) => {
        const items = prev[curr.result.spanId] || [];
        items.push(curr);

        return {
          ...prev,
          [curr.result.spanId]: items,
        };
      }, {});
  },

  getStructuredAssertion(assertion: string): TStructuredAssertion {
    const [left, right] = assertion.split(OperatorRegexp);
    const [comparator = CompareOperatorSymbolMap.EQUALS] = assertion.match(OperatorRegexp) ?? [];

    return {
      left,
      comparator: comparator as TCompareOperatorSymbol,
      right,
    };
  },

  getStringAssertion({left, comparator, right}: TStructuredAssertion): string {
    return `${left} ${comparator} ${right}`;
  },
});

export default AssertionService();
