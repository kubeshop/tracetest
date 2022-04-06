import {search} from 'jmespath';
import {SemanticGroupNames, SemanticGroupsSignature} from '../lib/SelectorDefaultAttributes';
import {ISpanAttributes, LOCATION_NAME, ResourceSpan} from '../types';

const SemanticGroupNamesList = Object.values(SemanticGroupNames);

export const getSpanValue = (
  span: ResourceSpan,
  locationName: LOCATION_NAME,
  valueType: string,
  key: string
): string => {
  switch (locationName) {
    case LOCATION_NAME.INSTRUMENTATION_LIBRARY:
    case LOCATION_NAME.RESOURCE_ATTRIBUTES: {
      const attributeList: ISpanAttributes = search(span, `resource.attributes[? key==\`${key}\`]`);

      return attributeList.find(attribute => attribute.key === key)?.value[valueType];
    }
    case LOCATION_NAME.SPAN:
    case LOCATION_NAME.SPAN_ATTRIBUTES: {
      const attributeList: ISpanAttributes = search(span, `instrumentationLibrarySpans[].spans[].attributes | []`);

      return attributeList.find(attribute => attribute.key === key)?.value[valueType];
    }
    // case LOCATION_NAME.SPAN_ID: {
    //   return `instrumentationLibrarySpans[?spans[?spanId == \`${spanId}\`]
    // }
    default:
      return '';
  }
};

export const getSpanType = (span: ResourceSpan) => {
  const [{spans: [{attributes = []}] = []} = {}] = span.instrumentationLibrarySpans || [];

  const findAttribute = (groupName: SemanticGroupNames) => attributes.find(({key}) => key.trim().startsWith(groupName));

  return SemanticGroupNamesList.find(groupName => Boolean(findAttribute(groupName))) || SemanticGroupNames.Http;
};

const getAttributeValueList = (span: ResourceSpan, attributeList: string[]) =>
  attributeList.reduce<string[]>((list, attribute) => {
    const value = getSpanValue(span, LOCATION_NAME.RESOURCE_ATTRIBUTES, 'stringValue', attribute);

    return value ? list.concat([value]) : list;
  }, []);

export const getSpanSignature = (span: ResourceSpan): string[] => {
  const type = getSpanType(span);

  const {SPAN_ATTRIBUTES, RESOURCE_ATTRIBUTES} = SemanticGroupsSignature[type];
  const spanAttributeList = getAttributeValueList(span, SPAN_ATTRIBUTES);
  const resourceAttributeList = getAttributeValueList(span, RESOURCE_ATTRIBUTES);
  const spanId = `#${span.instrumentationLibrarySpans[0]?.spans[0]?.spanId.slice(-4)}`;

  return [...resourceAttributeList, ...spanAttributeList, spanId];
};
