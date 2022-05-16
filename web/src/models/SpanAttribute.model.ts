import {isEmpty} from 'lodash';
import {SpanAttributeType} from '../constants/SpanAttribute.constants';
import {IRawSpanAttribute, ISpanAttribute} from '../types/SpanAttribute.types';

const spanAttributeTypeList = Object.values(SpanAttributeType);

const getSpanAttributeValueType = (attribute: IRawSpanAttribute): SpanAttributeType =>
  spanAttributeTypeList.find(type => {
    const value = attribute.value[type];
    if (typeof value === 'number') return true;
    return !isEmpty(value);
  }) || SpanAttributeType.stringValue;

const getSpanAttributeValue = (attribute: IRawSpanAttribute): string => {
  const attributeType = getSpanAttributeValueType(attribute);
  const value = attribute.value[attributeType];

  if (!value) return '<Empty value>';
  switch (attributeType) {
    case SpanAttributeType.kvlistValue: {
      return JSON.stringify(value);
    }

    default: {
      return String(value);
    }
  }
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
