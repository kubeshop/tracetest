import {ISpanAttribute} from './SpanAttribute.types';

export interface IInstrumentationLibrary {
  name: string;
  version: string;
}

export interface IResource {
  attributes: ISpanAttribute[];
}

export interface ISpanFlatAttribute {
  key: string;
  value: string;
}

export interface IInstrumentationLibrarySpan {
  instrumentationLibrary: IInstrumentationLibrary;
  spans: ISpan[];
}

export interface IResourceSpan {
  resource: IResource;
  instrumentationLibrarySpans: IInstrumentationLibrarySpan[];
}

export interface ISpan {
  traceId: string;
  spanId: string;
  parentSpanId: string;
  name: string;
  kind: number;
  startTimeUnixNano: string;
  endTimeUnixNano: string;
  attributes: ISpanAttribute[];
  status: {code: string};
  events: Event[];
}
