import {isEmpty} from 'lodash';
import {SpanAttributeType} from '../constants/SpanAttribute.constants';
import {IRawSpanAttribute, ISpanAttribute} from '../types/SpanAttribute.types';

const spanAttributeTypeList = Object.values(SpanAttributeType);

const typeOfList = ['number', 'boolean', 'string'];

const getSpanAttributeValueType = (attribute: IRawSpanAttribute): SpanAttributeType =>
  spanAttributeTypeList.find(type => {
    const value = attribute.value[type];

    return typeOfList.includes(typeof value) || !isEmpty(value);
  }) || SpanAttributeType.stringValue;

const getSpanAttributeValue = (attribute: IRawSpanAttribute): string => {
  const type = getSpanAttributeValueType(attribute);
  const value = attribute.value[type];

  if (['number', 'boolean'].includes(typeof value)) return String(value);

  return (
    (type === SpanAttributeType.kvlistValue && JSON.stringify(value)) || (value && String(value)) || '<Empty value>'
  );
};

const SpanAttribute = (rawAttribute: IRawSpanAttribute): ISpanAttribute => {
  const type = getSpanAttributeValueType(rawAttribute);
  const value = getSpanAttributeValue(rawAttribute);

  return {
    type,
    name: rawAttribute.key,
    value,
  };
};

export default SpanAttribute;
