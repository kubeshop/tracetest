import { TTraceSchemas } from 'types/Common.types';
import Span from './Span.model';

export type TRawTrace = TTraceSchemas['Trace'];
type Trace = {
  spans: Span[];
  traceId: string;
};

// TODO: keep the flat map of spans for easy access
const Trace = ({traceId = '', flat = {}}: TRawTrace): Trace => {
  return {
    traceId,
    spans: Object.values(flat).map(rawSpan => Span(rawSpan)),
  };
};

export default Trace;
