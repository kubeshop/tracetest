import {SpanAttributeType} from 'constants/SpanAttribute.constants';

export type IRawSpanAttributeValue = {
  stringValue?: string;
  intValue?: number;
  booleanValue?: boolean;
  doubleValue?: number;
  kvlistValue: {values: IRawSpanAttribute[]};
};

export type IRawSpanAttribute = {
  key: string;
  value: IRawSpanAttributeValue;
};

export type ISpanAttribute = {
  type: SpanAttributeType;
  name: string;
  value: string;
};
