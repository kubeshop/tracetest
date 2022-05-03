import {SemanticGroupNames} from '../constants/SemanticGroupNames.constants';
import {TItemSelector} from './Assertion.types';
import {Modify, TraceSchemas} from './Common.types';
import {TRawSpanAttribute, TSpanAttribute} from './SpanAttribute.types';

export type TInstrumentationLibrary = TraceSchemas['InstrumentationLibrary'];

export type TResource = Modify<
  TraceSchemas['Resource'],
  {
    attributes: TRawSpanAttribute[];
  }
>;

export interface TSpanFlatAttribute {
  key: string;
  value: string;
}

export type TInstrumentationLibrarySpan = TraceSchemas['InstrumentationLibrarySpans'];

export type TResourceSpan = TraceSchemas['ResourceSpans'];

export type TRawSpan = TraceSchemas['Span'];

export type TSpan = Modify<
  TRawSpan,
  {
    spanId: string;
    attributes: Record<string, TSpanAttribute>;
    instrumentationLibrary: TInstrumentationLibrary;
    type: SemanticGroupNames;
    duration: number;
    signature: TItemSelector[];
    attributeList: TSpanFlatAttribute[];
  }
>;
