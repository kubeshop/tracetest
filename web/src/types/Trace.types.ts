import {TraceSchemas} from './Common.types';
import {TSpan} from './Span.types';

export type TRawTrace = TraceSchemas['Trace'];

export type TTrace = {
  spans: TSpan[];
  description: string;
};
