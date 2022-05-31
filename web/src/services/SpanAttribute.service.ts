import {isEmpty, remove} from 'lodash';
import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {
  SectionNames,
  SelectorAttributesBlackList,
  SelectorAttributesWhiteList,
  SpanAttributeSections,
} from 'constants/Span.constants';
import {Attributes} from 'constants/SpanAttribute.constants';
import {TSpanFlatAttribute} from 'types/Span.types';
import {isJson} from 'utils/Common';

const flatAttributes = Object.values(Attributes);

const filterAttributeList = (attributeList: TSpanFlatAttribute[], attrKeyList: string[]): TSpanFlatAttribute[] => {
  return attrKeyList.reduce<TSpanFlatAttribute[]>((list, key) => {
    const foundAttrList = attributeList.filter(attr => attr.key === key);

    return foundAttrList.length ? list.concat(foundAttrList) : list;
  }, []);
};

const removeFromAttributeList = (attributeList: TSpanFlatAttribute[], attrKeyList: string[]): TSpanFlatAttribute[] =>
  remove(attributeList, attr => !attrKeyList.includes(attr.key));

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

  getFilteredSelectorAttributeList(
    attributeList: TSpanFlatAttribute[],
    currentSelectorList: string[]
  ): TSpanFlatAttribute[] {
    const duplicatedFiltered = removeFromAttributeList(attributeList, currentSelectorList);
    const whiteListFiltered = filterAttributeList(duplicatedFiltered, SelectorAttributesWhiteList);
    const blackListFiltered = removeFromAttributeList(whiteListFiltered, SelectorAttributesBlackList);
    const customList = attributeList.filter(attr => !flatAttributes.includes(attr.key));

    const parsedList = blackListFiltered.concat(customList).filter(attr => !isJson(attr.value) && !isEmpty(attr));

    return parsedList;
  },
});

export default SpanAttributeService();
