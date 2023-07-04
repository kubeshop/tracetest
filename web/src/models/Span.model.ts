import {
  SemanticGroupNames,
  SemanticGroupNamesToSystem,
  SemanticGroupsSignature,
} from 'constants/SemanticGroupNames.constants';
import {Attributes} from 'constants/SpanAttribute.constants';
import {SpanKind} from 'constants/Span.constants';
import {TSpanFlatAttribute} from 'types/Span.types';
import SpanAttribute from './SpanAttribute.model';
import {Model, TTraceSchemas} from '../types/Common.types';

const SemanticGroupNamesList = Object.values(SemanticGroupNames);
const SpanKindList = Object.values(SpanKind);

export type TRawSpan = TTraceSchemas['Span'];
type Span = Model<
  TRawSpan,
  {
    attributes: Record<string, SpanAttribute>;
    type: SemanticGroupNames;
    duration: string;
    signature: TSpanFlatAttribute[];
    attributeList: TSpanFlatAttribute[];
    children?: Span[];
    kind: SpanKind;
    service: string;
    system: string;
  }
>;

const getSpanSignature = (
  attributes: Record<string, SpanAttribute>,
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

const Span = ({id = '', attributes = {}, startTime = 0, endTime = 0, parentId = '', name = ''}: TRawSpan): Span => {
  const mappedAttributeList: TSpanFlatAttribute[] = [{key: 'name', value: name}];
  const attributeList = Object.entries(attributes)
    .map<TSpanFlatAttribute>(([key, value]) => ({
      value: String(value),
      key,
    }))
    .concat(mappedAttributeList);

  const attributesMap = attributeList.reduce((map, rawSpanAttribute) => {
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    return {...map, [spanAttribute.name]: SpanAttribute(rawSpanAttribute)};
  }, {});

  const kind = SpanKindList.find(spanKind => spanKind === attributes[Attributes.KIND]) || SpanKind.INTERNAL;
  const duration = attributes[Attributes.TRACETEST_SPAN_DURATION] || '0ns';
  const type =
    SemanticGroupNamesList.find(
      semanticGroupName => semanticGroupName === attributes[Attributes.TRACETEST_SPAN_TYPE]
    ) || SemanticGroupNames.General;
  const service = attributes?.[Attributes.SERVICE_NAME] ?? '';
  const system = attributes?.[SemanticGroupNamesToSystem[type]] ?? '';

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
    service,
    system,
  };
};

export default Span;
