import {TTraceSchemas} from './Common.types';
import {TSpan} from './Span.types';

export type TRawTrace = TTraceSchemas['Trace'];

export type TTrace = {
  spans: TSpan[];
  traceId: string;
};
