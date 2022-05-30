import {SemanticGroupNames} from '../constants/SemanticGroupNames.constants';
import {SectionNames, SpanAttributeSections} from '../constants/Span.constants';
import {Attributes} from '../constants/SpanAttribute.constants';
import {TSpanFlatAttribute} from '../types/Span.types';

const flatAttributes = Object.values(Attributes);

const filterAttributeList = (attributeList: TSpanFlatAttribute[], attrKeyList: string[]): TSpanFlatAttribute[] => {
  return attrKeyList.reduce<TSpanFlatAttribute[]>((list, key) => {
    const foundAttr = attributeList.find(attr => attr.key === key);

    return foundAttr ? list.concat(foundAttr) : list;
  }, []);
};

const SpanAttributeService = () => ({
  getSpanAttributeSectionsList(
    attributeList: TSpanFlatAttribute[],
    type: SemanticGroupNames
  ): {section: string; attributeList: TSpanFlatAttribute[]}[] {
    const sections = SpanAttributeSections[type] || {};
    const defaultSectionList = [
      {
        section: SectionNames.custom,
        attributeList: attributeList.filter(attr => !flatAttributes.includes(attr.key)),
      },
      {
        section: SectionNames.all,
        attributeList,
      },
    ];

    const sectionList = Object.entries(sections).reduce<{section: string; attributeList: TSpanFlatAttribute[]}[]>(
      (list, [key, attrKeyList]) =>
        list.concat([{section: key, attributeList: filterAttributeList(attributeList, attrKeyList)}]),
      []
    );

    return sectionList.concat(defaultSectionList);
  },
});

export default SpanAttributeService();
