import {SemanticGroupNames, SemanticGroupsSignature} from '../constants/SemanticGroupNames.constants';
import {LOCATION_NAME} from '../constants/Span.constants';
import {IItemSelector} from '../types/Assertion.types';
import {IInstrumentationLibrary, IRawSpan, IResource, IResourceSpan, ISpan} from '../types/Span.types';
import {ISpanAttribute} from '../types/SpanAttribute.types';
import SpanAttribute from './SpanAttribute.model';

const SemanticGroupNamesList = Object.values(SemanticGroupNames);

const getSpanType = ({attributes}: IRawSpan) => {
  const findAttribute = (groupName: SemanticGroupNames) => attributes.find(({key}) => key.trim().startsWith(groupName));

  return SemanticGroupNamesList.find(groupName => Boolean(findAttribute(groupName))) || SemanticGroupNames.General;
};

const getSpanSignature = (attributes: Record<string, ISpanAttribute>, type: SemanticGroupNames): IItemSelector[] => {
  const attributeNameList = SemanticGroupsSignature[type];

  return attributeNameList.reduce<IItemSelector[]>((list, attributeName) => {
    const attribute = attributes[attributeName];

    return attribute
      ? list.concat([
          {
            propertyName: attributeName,
            value: attribute.value,
            valueType: attribute.type,
            locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
          },
        ])
      : list;
  }, []);
};

const Span = (
  rawSpan: IRawSpan,
  {attributes: resourceAttributes = []}: IResource,
  instrumentationLibrary: IInstrumentationLibrary
): ISpan => {
  const attributesMap = rawSpan.attributes
    .concat(resourceAttributes)
    .reduce<Record<string, ISpanAttribute>>((map, rawSpanAttribute) => {
      const spanAttribute = SpanAttribute(rawSpanAttribute);

      return {...map, [spanAttribute.name]: SpanAttribute(rawSpanAttribute)};
    }, {});

  const duration = Number(
    ((Number(rawSpan.endTimeUnixNano) - Number(rawSpan.startTimeUnixNano)) / 1000 / 1000).toFixed(1)
  );
  const type = getSpanType(rawSpan);

  return {
    ...rawSpan,
    attributes: attributesMap,
    attributeList: Object.entries(attributesMap).map(([key, {value}]) => ({key, value})),
    instrumentationLibrary,
    type,
    duration,
    signature: getSpanSignature(attributesMap, type),
  };
};

Span.createFromResourceSpanList = (resourceSpans: IResourceSpan[]): ISpan[] =>
  resourceSpans.reduce<ISpan[]>((spanList, {resource, instrumentationLibrarySpans}) => {
    const spans = instrumentationLibrarySpans
      .flatMap<ISpan[]>(({instrumentationLibrary, spans: innerSpans}) =>
        innerSpans.map(span => Span(span, resource, instrumentationLibrary))
      )
      .flat();

    return spanList.concat(spans);
  }, []);

export default Span;
