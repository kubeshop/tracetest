import {SemanticGroupNames, SemanticGroupsSignature} from '../constants/SemanticGroupNames.constants';
import {TItemSelector} from '../types/Assertion.types';
import {TInstrumentationLibrary, TRawSpan, TResource, TResourceSpan, TSpan} from '../types/Span.types';
import {TSpanAttribute} from '../types/SpanAttribute.types';
import SpanAttribute from './SpanAttribute.model';

const SemanticGroupNamesList = Object.values(SemanticGroupNames);

const getSpanType = ({attributes = []}: TRawSpan) => {
  const findAttribute = (groupName: SemanticGroupNames) =>
    attributes.find(({key = ''}) => key.trim().startsWith(groupName));

  return SemanticGroupNamesList.find(groupName => Boolean(findAttribute(groupName))) || SemanticGroupNames.General;
};

const getSpanSignature = (attributes: Record<string, TSpanAttribute>, type: SemanticGroupNames): TItemSelector[] => {
  const attributeNameList = SemanticGroupsSignature[type];

  return attributeNameList.reduce<TItemSelector[]>((list, attributeName) => {
    const attribute = attributes[attributeName];

    return attribute
      ? list.concat([
          {
            propertyName: attributeName,
            value: attribute.value,
            valueType: attribute.type,
            locationName: 'SPAN_ATTRIBUTES',
          },
        ])
      : list;
  }, []);
};

const Span = (
  rawSpan: TRawSpan,
  {attributes: resourceAttributes = []}: TResource,
  instrumentationLibrary: TInstrumentationLibrary
): TSpan => {
  const {attributes = [], spanId = ''} = rawSpan;
  const attributesMap =
    attributes.concat(resourceAttributes).reduce<Record<string, TSpanAttribute>>((map, rawSpanAttribute) => {
      const spanAttribute = SpanAttribute(rawSpanAttribute);

      return {...map, [spanAttribute.name]: SpanAttribute(rawSpanAttribute)};
    }, {}) || {};

  const duration = (Number(rawSpan.endTimeUnixNano) - Number(rawSpan.startTimeUnixNano)) * 1000 * 1000;
  const type = getSpanType(rawSpan);

  return {
    ...rawSpan,
    spanId,
    attributes: attributesMap,
    attributeList: Object.entries(attributesMap).map(([key, {value}]) => ({key, value})),
    instrumentationLibrary,
    type,
    duration,
    signature: getSpanSignature(attributesMap, type),
  };
};

Span.createFromResourceSpanList = (resourceSpans: TResourceSpan[]): TSpan[] =>
  resourceSpans.reduce<TSpan[]>((spanList, {resource, instrumentationLibrarySpans = []}) => {
    const spans = instrumentationLibrarySpans
      .flatMap<TSpan[]>(({instrumentationLibrary, spans: innerSpans = []}) =>
        innerSpans.map(span => Span(span, resource, instrumentationLibrary!))
      )
      .flat();

    return spanList.concat(spans);
  }, []);

export default Span;
