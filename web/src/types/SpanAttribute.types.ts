import {Schemas, TraceSchemas} from './Common.types';

export type TRawSpanAttribute = TraceSchemas['Attribute'];

export type TSpanAttributeValueType = Schemas['SpanAssertion']['valueType']

export type TSpanAttribute = {
  type: Schemas['SpanAssertion']['valueType'];
  name: string;
  value: string;
};
