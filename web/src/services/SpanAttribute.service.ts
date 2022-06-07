import {isEmpty, remove} from 'lodash';
import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {
  SectionNames,
  SelectorAttributesBlackList,
  SelectorAttributesWhiteList,
  SpanAttributeSections,
} from 'constants/Span.constants';
import {Attributes, TraceTestAttributes} from 'constants/SpanAttribute.constants';
import {TSpanFlatAttribute} from 'types/Span.types';
import {isJson} from 'utils/Common';

const flatAttributes = Object.values(Attributes);
const flatTraceTestAttributes = Object.values(TraceTestAttributes);

const filterAttributeList = (attributeList: TSpanFlatAttribute[], attrKeyList: string[]): TSpanFlatAttribute[] => {
  return attrKeyList.reduce<TSpanFlatAttribute[]>((list, key) => {
    const foundAttrList = attributeList.filter(attr => attr.key.indexOf(key) === 0);

    return foundAttrList.length ? list.concat(foundAttrList) : list;
  }, []);
};

const removeFromAttributeList = (attributeList: TSpanFlatAttribute[], attrKeyList: string[]): TSpanFlatAttribute[] =>
  remove(attributeList, attr => !attrKeyList.includes(attr.key));

const getCustomAttributeList = (attributeList: TSpanFlatAttribute[]) => {
  const traceTestList = filterAttributeList(attributeList, flatTraceTestAttributes);

  const filetedList = attributeList.filter(attr => {
    const foundAttr = flatAttributes.find(key => attr.key.indexOf(key) === 0);

    return !foundAttr;
  });

  return traceTestList.concat(filetedList);
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
        attributeList: getCustomAttributeList(attributeList),
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
    const customList = getCustomAttributeList(attributeList);

    const parsedList = blackListFiltered.concat(customList).filter(attr => !isJson(attr.value) && !isEmpty(attr));

    return parsedList;
  },
});

export default SpanAttributeService();
