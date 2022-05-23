import {differenceBy, intersectionBy} from 'lodash';
import {SELECTOR_DEFAULT_ATTRIBUTES, SemanticGroupNameNodeMap} from '../constants/SemanticGroupNames.constants';
import {TSpan, TSpanFlatAttribute} from '../types/Span.types';

const itemSelectorKeys = SELECTOR_DEFAULT_ATTRIBUTES.flatMap(el => el.attributes);

const SpanService = () => ({
  getSpanNodeInfo(span: TSpan) {
    const signatureObject = span.signature.reduce<Record<string, string>>(
      (signature, {key, value}) => ({
        ...signature,
        [key]: value,
      }),
      {}
    );

    const {primary, type} = SemanticGroupNameNodeMap[span.type];

    const attributeKey = primary.find(key => signatureObject[key]) || '';

    return {
      primary: signatureObject[attributeKey] || '',
      heading: signatureObject[type] || '',
    };
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
});

export default SpanService();
