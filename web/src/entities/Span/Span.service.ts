import {search} from 'jmespath';
import {uniqBy} from 'lodash';
import {SemanticGroupNameNodeMap, SemanticGroupNames, SemanticGroupsSignature} from '../../constants/SemanticGroupNames.constants';
import {escapeString} from '../../utils/Common';
import { TItemSelector } from '../Assertion/Assertion.types';
import { TSpanAttribute } from '../SpanAttribute/SpanAttribute.types';
import {getSpanAttributeValue, getSpanAttributeValueType} from '../SpanAttribute/SpansAttribute.service';
import { TTrace } from '../Trace/Trace.types';
import { LOCATION_NAME } from './Span.constants';
import { TResourceSpan, TSpan, TSpanFlatAttribute } from './Span.types';

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

export const getSpanValue = (resourceSpan: TResourceSpan, locationName: LOCATION_NAME, key: string) => {
  const searchString = locationNameMap[locationName] || '';

  if (!searchString) return '';

  const attributeList: TSpanAttribute[] = search(resourceSpan, escapeString(searchString));
  const attributeFound = attributeList?.find(attribute => attribute.key === key);

  return attributeFound ? getSpanAttributeValue(attributeFound) : '';
};

export const getSpanAttributeType = (resourceSpan: TResourceSpan, locationName: LOCATION_NAME, key: string) => {
  const searchString = locationNameMap[locationName] || '';

  if (!searchString) return '';

  const attributeList: TSpanAttribute[] = search(resourceSpan, escapeString(searchString));
  const attributeFound = attributeList?.find(attribute => attribute.key === key);

  return attributeFound ? getSpanAttributeValueType(attributeFound) : '';
};

export const getResourceSpanBySpanId = (spanId: string, trace: TTrace): TResourceSpan | undefined => {
  const [resourceSpan]: Array<TResourceSpan> = search(
    trace,
    `resourceSpans|[]| [? instrumentationLibrarySpans[?spans[?spanId == \`${spanId}\`]]]`
  );

  return resourceSpan;
};

export const getSpanType = (spanId: string, trace: TTrace) => {
  const resourceSpan = getResourceSpanBySpanId(spanId, trace);
  if (!resourceSpan) return undefined;

  const {attributes = []} =
    resourceSpan.instrumentationLibrarySpans.flatMap(({spans}) => spans).find(({spanId: id}) => id === spanId) || {};

  const findAttribute = (groupName: SemanticGroupNames) => attributes.find(({key}) => key.trim().startsWith(groupName));

  return SemanticGroupNamesList.find(groupName => Boolean(findAttribute(groupName))) || SemanticGroupNames.General;
};

const getAttributeValueList = (resourceSpan: TResourceSpan, attributeList: string[], locationName: LOCATION_NAME) =>
  attributeList.reduce<TItemSelector[]>((list, attribute) => {
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

export const getSpan = (resourceSpan: TResourceSpan): TSpan | undefined =>
  resourceSpan.instrumentationLibrarySpans[0]?.spans[0];

export const getSpanSignature = (spanId: string, trace: TTrace): TItemSelector[] => {
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

export const getSpanNodeInfo = (spanId: string, trace: TTrace) => {
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

export const getSpanAttributeList = (resourceSpan: TResourceSpan): TSpanFlatAttribute[] => {
  const spanAttributeList: TSpanAttribute[] = search(resourceSpan, escapeString(spanAttributeSearch)) || [];
  const resourceAttributeList: TSpanAttribute[] = search(resourceSpan, escapeString(resourceAttributeSearch)) || [];
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
