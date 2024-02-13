import {differenceBy, intersectionBy} from 'lodash';
import {CompareOperator} from 'constants/Operator.constants';
import {SELECTOR_DEFAULT_ATTRIBUTES, SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {SpanKind} from 'constants/Span.constants';
import Span from 'models/Span.model';
import {TSpanFlatAttribute} from 'types/Span.types';
import {getObjectIncludesText} from 'utils/Common';
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

  // TODO: this is very costly, we might need to move this to the backend
  searchSpanList(spanList: Span[], searchText: string) {
    if (!searchText.trim()) return [];

    return spanList.reduce<string[]>(
      (matchList, span) => (getObjectIncludesText(span.attributes, searchText) ? [...matchList, span.id] : matchList),
      []
    );
  },
});

export default SpanService();
