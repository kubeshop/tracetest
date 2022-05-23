import {differenceInSeconds, parseISO} from 'date-fns';
import {SemanticGroupNames, SemanticGroupsSignature} from '../constants/SemanticGroupNames.constants';

import {TRawSpan, TSpan, TSpanFlatAttribute} from '../types/Span.types';
import {TSpanAttribute} from '../types/SpanAttribute.types';
import SpanAttribute from './SpanAttribute.model';

const SemanticGroupNamesList = Object.values(SemanticGroupNames);

const getSpanType = (attributeList: TSpanFlatAttribute[]) => {
  const findAttribute = (groupName: SemanticGroupNames) =>
    attributeList.find(({key}) => key.trim().startsWith(groupName));

  return SemanticGroupNamesList.find(groupName => Boolean(findAttribute(groupName))) || SemanticGroupNames.General;
};

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

const Span = ({id = '', name = '', attributes = {}, startTime = '', endTime = '', parentId = ''}: TRawSpan): TSpan => {
  const attributeList = Object.entries(attributes).map<TSpanFlatAttribute>(([key, value]) => ({
    value: String(value),
    key,
  }));

  const attributesMap = attributeList.reduce<Record<string, TSpanAttribute>>((map, rawSpanAttribute) => {
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    return {...map, [spanAttribute.name]: SpanAttribute(rawSpanAttribute)};
  }, {});

  const duration = startTime && endTime ? differenceInSeconds(parseISO(startTime), parseISO(endTime)) + 1 : 0;

  const type = getSpanType(attributeList);

  return {
    id,
    parentId,
    name,
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
