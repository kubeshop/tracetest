import {
  OtelReference,
  OtelReferenceModel,
} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import AttributesTags from 'constants/AttributesTags.json';
import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {
  SectionNames,
  SelectorAttributesBlackList,
  SelectorAttributesWhiteList,
  SpanAttributeSections,
} from 'constants/Span.constants';
import {Attributes, TraceTestAttributes} from 'constants/SpanAttribute.constants';
import {isEmpty, remove} from 'lodash';
import {TSpanFlatAttribute} from 'types/Span.types';
import {isJson} from 'utils/Common';

const attributesTags: Record<string, OtelReferenceModel> = AttributesTags;
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
  getPseudoAttributeList: (count: number): TSpanFlatAttribute[] => [
    {key: TraceTestAttributes.TRACETEST_SELECTED_SPANS_COUNT, value: count.toString()},
  ],

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

    return blackListFiltered.concat(customList).filter(attr => !isJson(attr.value) && !isEmpty(attr));
  },
  referencePicker(reference: OtelReference, key: string): OtelReferenceModel {
    return reference[key] || attributesTags[key] || {description: '', tags: []};
  },
});

export default SpanAttributeService();
