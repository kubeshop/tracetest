import {search} from 'jmespath';
import {SemanticGroupNameNodeMap, SemanticGroupNames, SemanticGroupsSignature} from '../lib/SelectorDefaultAttributes';
import {ISpanAttribute, ItemSelector, ITrace, LOCATION_NAME, ResourceSpan, Span, TSpanAttributesList} from '../types';
import {escapeString} from '../utils';
import {getSpanAttributeValue, getSpanAttributeValueType} from './SpanAttributeService';

const SemanticGroupNamesList = Object.values(SemanticGroupNames);

const spanAttributeSearch = 'instrumentationLibrarySpans[].spans[].attributes | []';
const resourceAttributeSearch = 'resource.attributes | []';

const locationNameMap = {
  [LOCATION_NAME.INSTRUMENTATION_LIBRARY]: resourceAttributeSearch,
  [LOCATION_NAME.RESOURCE_ATTRIBUTES]: resourceAttributeSearch,
  [LOCATION_NAME.SPAN]: spanAttributeSearch,
  [LOCATION_NAME.SPAN_ATTRIBUTES]: spanAttributeSearch,
  [LOCATION_NAME.SPAN_ID]: '',
};

export const getSpanValue = (resourceSpan: ResourceSpan, locationName: LOCATION_NAME, key: string) => {
  const searchString = locationNameMap[locationName] || '';

  if (!searchString) return '';

  const attributeList: ISpanAttribute[] = search(resourceSpan, escapeString(searchString));
  const attributeFound = attributeList?.find(attribute => attribute.key === key);

  return attributeFound ? getSpanAttributeValue(attributeFound) : '';
};

export const getSpanAttributeType = (resourceSpan: ResourceSpan, locationName: LOCATION_NAME, key: string) => {
  const searchString = locationNameMap[locationName] || '';

  if (!searchString) return '';

  const attributeList: ISpanAttribute[] = search(resourceSpan, escapeString(searchString));
  const attributeFound = attributeList?.find(attribute => attribute.key === key);

  return attributeFound ? getSpanAttributeValueType(attributeFound) : '';
};

export const getResourceSpanBySpanId = (spanId: string, trace: ITrace): ResourceSpan | undefined => {
  const [resourceSpan]: Array<ResourceSpan> = search(
    trace,
    `resourceSpans|[]| [? instrumentationLibrarySpans[?spans[?spanId == \`${spanId}\`]]]`
  );

  return resourceSpan;
};

export const getSpanType = (spanId: string, trace: ITrace) => {
  const resourceSpan = getResourceSpanBySpanId(spanId, trace);
  if (!resourceSpan) return undefined;

  const {attributes = []} =
    resourceSpan.instrumentationLibrarySpans.flatMap(({spans}) => spans).find(({spanId: id}) => id === spanId) || {};

  const findAttribute = (groupName: SemanticGroupNames) => attributes.find(({key}) => key.trim().startsWith(groupName));

  return SemanticGroupNamesList.find(groupName => Boolean(findAttribute(groupName))) || SemanticGroupNames.General;
};

const getAttributeValueList = (resourceSpan: ResourceSpan, attributeList: string[], locationName: LOCATION_NAME) =>
  attributeList.reduce<ItemSelector[]>((list, attribute) => {
    const value = getSpanValue(resourceSpan, locationName, attribute);
    const valueType = getSpanAttributeType(resourceSpan, locationName, attribute);

    return value
      ? list.concat([
          {
            propertyName: attribute,
            value,
            valueType,
            locationName,
          },
        ])
      : list;
  }, []);

export const getSpan = (resourceSpan: ResourceSpan): Span | undefined =>
  resourceSpan.instrumentationLibrarySpans[0]?.spans[0];

export const getSpanSignature = (spanId: string, trace: ITrace): ItemSelector[] => {
  const type = getSpanType(spanId, trace);
  const resourceSpan = getResourceSpanBySpanId(spanId, trace);

  if (!type || !resourceSpan) return [];

  const {SPAN_ATTRIBUTES: spanAttributes, RESOURCE_ATTRIBUTES: resourceAttributes} = SemanticGroupsSignature[type];
  const spanAttributeList = getAttributeValueList(resourceSpan, spanAttributes, LOCATION_NAME.SPAN_ATTRIBUTES);
  const resourceAttributeList = getAttributeValueList(
    resourceSpan,
    resourceAttributes,
    LOCATION_NAME.RESOURCE_ATTRIBUTES
  );

  return [...resourceAttributeList, ...spanAttributeList];
};

export const getSpanNodeInfo = (spanId: string, trace: ITrace) => {
  const spanType = getSpanType(spanId, trace);
  const signatureArray = getSpanSignature(spanId, trace);

  const signatureObject = signatureArray.reduce<Record<string, string>>(
    (signature, {propertyName, value}) => ({
      ...signature,
      [propertyName]: value,
    }),
    {}
  );

  const {primary, type} = SemanticGroupNameNodeMap[spanType!];

  const attributeKey = primary.find(key => signatureObject[key]) || '';

  return {
    primary: signatureObject[attributeKey] || '',
    heading: signatureObject[type],
    spanType,
  };
};

export const getSpanAttributeList = (resourceSpan: ResourceSpan): TSpanAttributesList => {
  const spanAttributeList: ISpanAttribute[] = search(resourceSpan, escapeString(spanAttributeSearch)) || [];
  const resourceAttributeList: ISpanAttribute[] = search(resourceSpan, escapeString(resourceAttributeSearch)) || [];
  const {spanId, parentSpanId, traceId, kind, status, name} = getSpan(resourceSpan) || {};

  const attributeList = [...resourceAttributeList, ...spanAttributeList].map(attribute => ({
    key: attribute.key,
    value: getSpanAttributeValue(attribute),
  }));

  const spanFieldList = Object.entries({traceId, spanId, parentSpanId, name, kind}).map(([key, value]) => ({
    key,
    value: typeof value !== 'undefined' ? String(value) : '',
  }));

  return [
    ...spanFieldList,
    {
      key: 'status',
      value: status?.code!,
    },
    ...attributeList,
  ];
};
