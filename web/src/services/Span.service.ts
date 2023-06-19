import {differenceBy, intersectionBy} from 'lodash';
import {CompareOperator} from 'constants/Operator.constants';
import {SELECTOR_DEFAULT_ATTRIBUTES, SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {SpanKind} from 'constants/Span.constants';
import {TSpanFlatAttribute} from 'types/Span.types';
import {getObjectIncludesText} from 'utils/Common';
import {TResultAssertions, TResultAssertionsSummary} from 'types/Assertion.types';
import LinterResult from 'models/LinterResult.model';
import Span from 'models/Span.model';
import OperatorService from './Operator.service';

const itemSelectorKeys = SELECTOR_DEFAULT_ATTRIBUTES.flatMap(el => el.attributes);

const SpanService = () => ({
  getSpanInfo(span?: Span) {
    const kind = span?.kind ?? SpanKind.INTERNAL;
    const name = span?.name ?? '';
    const service = span?.service ?? '';
    const system = span?.system ?? '';
    const type = span?.type ?? SemanticGroupNames.General;

    return {kind, name, service, system, type};
  },

  getSelectedSpanListAttributes({attributeList}: Span, selectedSpanList: Span[]) {
    const intersectedAttributeList = intersectionBy(...selectedSpanList.map(el => el.attributeList), 'key');

    const selectedSpanAttributeList = attributeList?.reduce<TSpanFlatAttribute[]>((acc, item) => {
      if (itemSelectorKeys.includes(item.key)) return acc;

      return acc.concat([item]);
    }, []);

    return {
      intersectedList: intersectedAttributeList,
      differenceList: differenceBy(selectedSpanAttributeList, intersectedAttributeList, 'key'),
    };
  },

  getSelectorInformation(span: Span) {
    return `span[${(
      span?.signature.reduce<string>(
        (selector, {value, key}) =>
          `${selector}${key}${OperatorService.getOperatorSymbol(CompareOperator.EQUALS)}"${value}" `,
        ''
      ) || ''
    ).trim()}]`;
  },

  searchSpanList(spanList: Span[], searchText: string) {
    if (!searchText.trim()) return [];

    return spanList.reduce<string[]>(
      (matchList, span) => (getObjectIncludesText(span.attributes, searchText) ? [...matchList, span.id] : matchList),
      []
    );
  },

  getAssertionResultSummary(assertions: TResultAssertions): TResultAssertionsSummary {
    const resultSummary = Object.values(assertions).reduce<TResultAssertionsSummary>(
      ({failed: prevFailed, passed: prevPassed}, {failed, passed}) => ({
        failed: prevFailed.concat(failed),
        passed: prevPassed.concat(passed),
      }),
      {
        failed: [],
        passed: [],
      }
    );

    return resultSummary;
  },

  getLintBySpan(linterResult: LinterResult): TLintBySpan {
    return linterResult.plugins
      .flatMap(plugin => plugin.rules.map(rule => ({...rule, pluginName: plugin.name})))
      .flatMap(rule => rule.results.map(result => ({...result, ruleName: rule.name, pluginName: rule.pluginName})))
      .reduce((prev: TLintBySpan, curr) => {
        const value = prev[curr.spanId] || [];
        return {...prev, [curr.spanId]: [...value, curr]};
      }, {});
  },

  filterLintErrorsBySpan(linterResultsBySpan: TLintBySpan, spanId: string) {
    const results = linterResultsBySpan[spanId];
    return results?.filter(result => !result.passed) ?? [];
  },
});

export type TLintBySpanContent = {
  ruleName: string;
  pluginName: string;
  passed: boolean;
  spanId: string;
  errors: {
    error?: string;
    value?: string;
    expected?: string;
    level?: string;
    description?: string;
    suggestions?: string[];
  }[];
  severity: 'error' | 'warning';
};

export type TLintBySpan = Record<string, TLintBySpanContent[]>;

export default SpanService();
