import {differenceBy, intersectionBy} from 'lodash';
import {CompareOperator, PseudoSelector} from 'constants/Operator.constants';
import {SELECTOR_DEFAULT_ATTRIBUTES, SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {SpanKind} from 'constants/Span.constants';
import {TSpan, TSpanFlatAttribute} from 'types/Span.types';
import {getObjectIncludesText} from 'utils/Common';
import OperatorService from './Operator.service';

const itemSelectorKeys = SELECTOR_DEFAULT_ATTRIBUTES.flatMap(el => el.attributes);

const SpanService = () => ({
  getSpanInfo(span?: TSpan) {
    const kind = span?.kind ?? SpanKind.INTERNAL;
    const name = span?.name ?? '';
    const service = span?.service ?? '';
    const system = span?.system ?? '';
    const type = span?.type ?? SemanticGroupNames.General;

    return {kind, name, service, system, type};
  },

  getSelectedSpanListAttributes({attributeList}: TSpan, selectedSpanList: TSpan[]) {
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

  getSelectorInformation(span: TSpan) {
    const selectorList =
      span?.signature.map(attribute => ({
        value: attribute.value,
        key: attribute.key,
        operator: OperatorService.getOperatorSymbol(CompareOperator.EQUALS),
      })) || [];

    const pseudoSelector = {
      selector: PseudoSelector.ALL,
    };

    return {selectorList, pseudoSelector};
  },

  searchSpanList(spanList: TSpan[], searchText: string) {
    if (!searchText.trim()) return [];

    return spanList.reduce<string[]>(
      (matchList, span) => (getObjectIncludesText(span.attributes, searchText) ? [...matchList, span.id] : matchList),
      []
    );
  },
});

export default SpanService();
