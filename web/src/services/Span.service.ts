import {search} from 'jmespath';
import {uniqBy} from 'lodash';
import {
  SemanticGroupNameNodeMap,
  SemanticGroupNames,
  SemanticGroupsSignature,
} from '../constants/SemanticGroupNames.constants';
import {escapeString} from '../utils/Common';
import {IItemSelector} from '../types/Assertion.types';
import {ISpanAttribute} from '../types/SpanAttribute.types';
import {getSpanAttributeValue, getSpanAttributeValueType} from './SpansAttribute.service';
import {ITrace} from '../types/Trace.types';
import {LOCATION_NAME} from '../constants/Span.constants';
import {IResourceSpan, ISpan, ISpanFlatAttribute} from '../types/Span.types';

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

export const getSpanValue = (resourceSpan: IResourceSpan, locationName: LOCATION_NAME, key: string) => {
  const searchString = locationNameMap[locationName] || '';

  if (!searchString) return '';

  const attributeList: ISpanAttribute[] = search(resourceSpan, escapeString(searchString));
  const attributeFound = attributeList?.find(attribute => attribute.key === key);

  return attributeFound ? getSpanAttributeValue(attributeFound) : '';
};

export const getSpanAttributeType = (resourceSpan: IResourceSpan, locationName: LOCATION_NAME, key: string) => {
  const searchString = locationNameMap[locationName] || '';

  if (!searchString) return '';

  const attributeList: ISpanAttribute[] = search(resourceSpan, escapeString(searchString));
  const attributeFound = attributeList?.find(attribute => attribute.key === key);

  return attributeFound ? getSpanAttributeValueType(attributeFound) : '';
};

export const getResourceSpanBySpanId = (spanId: string, trace: ITrace): IResourceSpan | undefined => {
  const [resourceSpan]: Array<IResourceSpan> = search(
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

const getAttributeValueList = (resourceSpan: IResourceSpan, attributeList: string[], locationName: LOCATION_NAME) =>
  attributeList.reduce<IItemSelector[]>((list, attribute) => {
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

export const getSpan = (resourceSpan: IResourceSpan): ISpan | undefined =>
  resourceSpan.instrumentationLibrarySpans[0]?.spans[0];

export const getSpanSignature = (spanId: string, trace: ITrace): IItemSelector[] => {
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

  return [...spanAttributeList, ...resourceAttributeList];
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

export const getSpanAttributeList = (resourceSpan: IResourceSpan): ISpanFlatAttribute[] => {
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

  return uniqBy(
    [
      ...spanFieldList,
      {
        key: 'status',
        value: status?.code!,
      },
      ...attributeList,
    ],
    'key'
  );
};
