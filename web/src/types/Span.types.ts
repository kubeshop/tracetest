import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {SpanKind} from 'constants/Span.constants';
import {Model, TTraceSchemas} from './Common.types';
import {TSpanAttribute} from './SpanAttribute.types';

export type TSpanFlatAttribute = {
  key: string;
  value: string;
};

export type TRawSpan = TTraceSchemas['Span'];

export type TSpan = Model<
  TRawSpan,
  {
    attributes: Record<string, TSpanAttribute>;
    type: SemanticGroupNames;
    duration: string;
    signature: TSpanFlatAttribute[];
    attributeList: TSpanFlatAttribute[];
    children?: TSpan[];
    kind: SpanKind;
    service: string;
    system: string;
  }
>;

export interface ISpanState {
  affectedSpans: string[];
  focusedSpan: string;
  selectedSpan?: TSpan;
  searchText: string;
  matchedSpans: string[];
}

export type TSpansResult = Record<
  string,
  {
    failed: number;
    passed: number;
  }
>;
