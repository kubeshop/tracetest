import {TSpanAttribute} from '../SpanAttribute/SpanAttribute.types';

export type TInstrumentationLibrary = {
  name: string;
  version: string;
};

export type TResource = {
  attributes: TSpanAttribute[];
};

export type TSpanFlatAttribute = {
  key: string,
  value: string
};

export type TInstrumentationLibrarySpan = {
  instrumentationLibrary: TInstrumentationLibrary;
  spans: TSpan[];
};

export type TResourceSpan = {
  resource: TResource;
  instrumentationLibrarySpans: TInstrumentationLibrarySpan[];
};

export type TSpan = {
  traceId: string;
  spanId: string;
  parentSpanId: string;
  name: string;
  kind: number;
  startTimeUnixNano: string;
  endTimeUnixNano: string;
  attributes: TSpanAttribute[];
  status: {code: string};
  events: Event[];
};
