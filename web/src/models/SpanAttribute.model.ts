import {isEmpty} from 'lodash';
import {TRawSpanAttribute, TSpanAttribute, TSpanAttributeValueType} from '../types/SpanAttribute.types';

const spanAttributeTypeList: TSpanAttributeValueType[] = ['boolValue', 'doubleValue', 'intValue', 'kvlistValue', 'stringValue'];

const getSpanAttributeValueType = (attribute: TRawSpanAttribute): TSpanAttributeValueType =>
  spanAttributeTypeList.find(type => {
    const value = attribute.value[type];
    if (typeof value === 'number') return true;
    return !isEmpty(value);
  }) || 'stringValue';

const getSpanAttributeValue = (attribute: TRawSpanAttribute): string => {
  const attributeType = getSpanAttributeValueType(attribute);
  const value = attribute.value[attributeType];

  if (!value) return 'Empty value';
  switch (attributeType) {
    case 'kvlistValue': {
      return JSON.stringify(value);
    }

    default: {
      return String(value);
    }
  }
};

const SpanAttribute = (rawAttribute: TRawSpanAttribute): TSpanAttribute => {
  const type = getSpanAttributeValueType(rawAttribute);
  const value = getSpanAttributeValue(rawAttribute);

  return {
    type,
    name: rawAttribute.key || '',
    value,
  };
};

export default SpanAttribute;
