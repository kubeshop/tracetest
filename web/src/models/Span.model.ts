import {SemanticGroupNames, SemanticGroupsSignature} from 'constants/SemanticGroupNames.constants';
import {Attributes} from 'constants/SpanAttribute.constants';
import {SpanKind} from 'constants/Span.constants';
import {TRawSpan, TSpan, TSpanFlatAttribute} from 'types/Span.types';
import {TSpanAttribute} from 'types/SpanAttribute.types';
import SpanAttribute from './SpanAttribute.model';

const SemanticGroupNamesList = Object.values(SemanticGroupNames);
const SpanKindList = Object.values(SpanKind);

const getSpanSignature = (
  attributes: Record<string, TSpanAttribute>,
  type: SemanticGroupNames
): TSpanFlatAttribute[] => {
  const attributeNameList = SemanticGroupsSignature[type];

  return attributeNameList.reduce<TSpanFlatAttribute[]>((list, attributeName) => {
    const attribute = attributes[attributeName];

    return attribute
      ? list.concat([
          {
            key: attributeName,
            value: attribute.value,
          },
        ])
      : list;
  }, []);
};

const Span = ({id = '', attributes = {}, startTime = 0, endTime = 0, parentId = ''}: TRawSpan): TSpan => {
  const attributeList = Object.entries(attributes).map<TSpanFlatAttribute>(([key, value]) => ({
    value: String(value),
    key,
  }));

  const attributesMap = attributeList.reduce<Record<string, TSpanAttribute>>((map, rawSpanAttribute) => {
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    return {...map, [spanAttribute.name]: SpanAttribute(rawSpanAttribute)};
  }, {});

  const name = attributes[Attributes.NAME] || '';
  const kind = SpanKindList.find(spanKind => spanKind === attributes[Attributes.KIND]) || SpanKind.INTERNAL;
  const duration = Number(attributes[Attributes.TRACETEST_SPAN_DURATION]) || 0;
  const type =
    SemanticGroupNamesList.find(
      semanticGroupName => semanticGroupName === attributes[Attributes.TRACETEST_SPAN_TYPE]
    ) || SemanticGroupNames.General;

  return {
    id,
    parentId,
    name,
    kind,
    startTime,
    endTime,
    attributes: attributesMap,
    attributeList,
    type,
    duration,
    signature: getSpanSignature(attributesMap, type),
  };
};

export default Span;
