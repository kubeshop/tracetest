import {
  OtelReference,
  OtelReferenceModel,
} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import {SelectorAttributesBlackList, SelectorAttributesWhiteList} from 'constants/Span.constants';
import {Attributes, TraceTestAttributes} from 'constants/SpanAttribute.constants';
import {isEmpty, remove} from 'lodash';
import {TSpanFlatAttribute} from 'types/Span.types';
import {getObjectIncludesText, isJson} from 'utils/Common';
import {TResultAssertions, TResultAssertionsSummary} from 'types/Assertion.types';
import TestOutput from 'models/TestOutput.model';

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

  getReferencePicker(reference: OtelReference, attrName: string): OtelReferenceModel {
    return reference[attrName] || {description: '', note: '', tags: []};
  },

  getItMatchesAttributeByKey(reference: OtelReference, attrName: string, search: string): boolean {
    const {tags = [], description} = reference[attrName] || {description: '', tags: []};

    const availableTagsMatchInput = Boolean(
      tags.find(tag => tag.toString().toLowerCase().includes(search.toLowerCase()))
    );
    const currentOptionMatchInput = attrName.includes(search);
    const currentDescriptionMatchInput = description.includes(search);

    return availableTagsMatchInput || currentOptionMatchInput || currentDescriptionMatchInput;
  },

  filterAttributes(attributes: TSpanFlatAttribute[], searchText: string, semanticConventions: OtelReference) {
    if (!searchText.length) return attributes;

    const searchTextLowerCase = searchText.toLowerCase();
    return attributes.filter(({key, value}) => {
      const {description = '', tags = []} = semanticConventions?.[key] ?? {};

      return (
        key.toLowerCase().includes(searchTextLowerCase) ||
        value.toLowerCase().includes(searchTextLowerCase) ||
        description.toLowerCase().includes(searchTextLowerCase) ||
        getObjectIncludesText(tags, searchTextLowerCase)
      );
    });
  },

  getAttributeAssertionResults(attrName: string, assertions: TResultAssertions): TResultAssertionsSummary {
    const resultList = Object.entries(assertions).reduce<TResultAssertionsSummary>(
      ({failed: prevFailed, passed: prevPassed}, [assertionString, {failed, passed}]) => {
        const itMatches = assertionString.toLowerCase().includes(attrName.toLowerCase());

        return itMatches
          ? {
              failed: prevFailed.concat(failed),
              passed: prevPassed.concat(passed),
            }
          : {failed: prevFailed, passed: prevPassed};
      },
      {
        failed: [],
        passed: [],
      }
    );

    return resultList;
  },

  getOutputsFromAttributeName(attributeName: string, outputs: TestOutput[]): TestOutput[] {
    return outputs.filter(({value}) => value.toLowerCase().includes(attributeName.toLowerCase()));
  },
});

export default SpanAttributeService();
