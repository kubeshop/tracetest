import {search} from 'jmespath';
import {SemanticGroupNames, SemanticGroupsSignature} from '../lib/SelectorDefaultAttributes';
import {
  ISpanAttribute,
  ISpanAttributeValue,
  ItemSelector,
  ITrace,
  LOCATION_NAME,
  ResourceSpan,
  Span,
  TSpanAttributesList,
} from '../types';

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

export const getSpanValue = (
  resourceSpan: ResourceSpan,
  locationName: LOCATION_NAME,
  valueType: keyof ISpanAttributeValue,
  key: string
): string => {
  const searchString = locationNameMap[locationName] || '';

  if (!searchString) return '';

  const attributeList: ISpanAttribute[] = search(resourceSpan, searchString);
  const {value} = attributeList?.find(attribute => attribute.key === key) || {};

  return value ? String(value[valueType]) : '';
};

export const getSpanType = (resourceSpan: ResourceSpan) => {
  const [{spans: [{attributes = []}] = []} = {}] = resourceSpan.instrumentationLibrarySpans || [];

  const findAttribute = (groupName: SemanticGroupNames) => attributes.find(({key}) => key.trim().startsWith(groupName));

  return SemanticGroupNamesList.find(groupName => Boolean(findAttribute(groupName))) || SemanticGroupNames.General;
};

const getAttributeValueList = (resourceSpan: ResourceSpan, attributeList: string[], locationName: LOCATION_NAME) =>
  attributeList.reduce<ItemSelector[]>((list, attribute) => {
    const value = getSpanValue(resourceSpan, locationName, 'stringValue', attribute);

    return value
      ? list.concat([
          {
            propertyName: attribute,
            value,
            valueType: 'stringValue',
            locationName,
          },
        ])
      : list;
  }, []);

export const getResourceSpanBySpanId = (spanId: string, trace: ITrace): ResourceSpan | undefined => {
  const [resourceSpan]: Array<ResourceSpan> = search(
    trace,
    `resourceSpans|[]| [? instrumentationLibrarySpans[?spans[?spanId == \`${spanId}\`]]]`
  );

  return resourceSpan;
};

export const getSpan = (resourceSpan: ResourceSpan): Span | undefined =>
  resourceSpan.instrumentationLibrarySpans[0]?.spans[0];

export const getSpanSignature = (spanId: string, trace: ITrace): ItemSelector[] => {
  const resourceSpan = getResourceSpanBySpanId(spanId, trace);
  if (!resourceSpan) return [];

  const type = getSpanType(resourceSpan);

  const {SPAN_ATTRIBUTES: spanAttributes, RESOURCE_ATTRIBUTES: resourceAttributes} = SemanticGroupsSignature[type];
  const spanAttributeList = getAttributeValueList(resourceSpan, spanAttributes, LOCATION_NAME.SPAN_ATTRIBUTES);
  const resourceAttributeList = getAttributeValueList(
    resourceSpan,
    resourceAttributes,
    LOCATION_NAME.RESOURCE_ATTRIBUTES
  );

  return [...resourceAttributeList, ...spanAttributeList];
};

export const getSpanAttributeList = (resourceSpan: ResourceSpan): TSpanAttributesList => {
  const spanAttributeList: ISpanAttribute[] = search(resourceSpan, spanAttributeSearch) || [];
  const resourceAttributeList: ISpanAttribute[] = search(resourceSpan, resourceAttributeSearch) || [];
  const {spanId, parentSpanId, traceId, kind, status, name} = getSpan(resourceSpan) || {};

  const attributeList = [...resourceAttributeList, ...spanAttributeList].map(
    ({key, value: {intValue, stringValue, booleanValue}}) => ({
      key,
      value: String(intValue || stringValue || booleanValue),
    })
  );

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
