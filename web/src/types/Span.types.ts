import {SemanticGroupNames} from '../constants/SemanticGroupNames.constants';
import {Modify} from './Common.types';
import {IItemSelector} from './Assertion.types';
import {IRawSpanAttribute, ISpanAttribute} from './SpanAttribute.types';

export interface IInstrumentationLibrary {
  name: string;
  version: string;
}

export interface IResource {
  attributes: IRawSpanAttribute[];
}

export interface ISpanFlatAttribute {
  key: string;
  value: string;
}

export interface IInstrumentationLibrarySpan {
  instrumentationLibrary: IInstrumentationLibrary;
  spans: IRawSpan[];
}

export interface IResourceSpan {
  resource: IResource;
  instrumentationLibrarySpans: IInstrumentationLibrarySpan[];
}

export interface IRawSpan {
  traceId: string;
  spanId: string;
  parentSpanId?: string;
  name: string;
  kind: string;
  startTimeUnixNano: string;
  endTimeUnixNano: string;
  attributes: IRawSpanAttribute[];
  status: {code: string};
}

export type ISpan = Modify<
  IRawSpan,
  {
    attributes: Record<string, ISpanAttribute>;
    type: SemanticGroupNames;
    duration: number;
    signature: IItemSelector[];
    attributeList: ISpanFlatAttribute[];
  }
>;
