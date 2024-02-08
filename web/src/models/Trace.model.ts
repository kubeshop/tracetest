import {TTraceSchemas} from 'types/Common.types';
import Span from './Span.model';

export type TRawTrace = TTraceSchemas['Trace'];
export type TSpanMap = Record<string, Span>;
type Trace = {
  flat: TSpanMap;
  spans: Span[];
  traceId: string;
  rootSpan: Span;
};

const defaultTrace: TRawTrace = {
  traceId: '',
  flat: {},
  tree: {},
};

const Trace = ({traceId = '', flat = {}, tree = {}} = defaultTrace): Trace => ({
  traceId,
  rootSpan: Span(tree),
  flat: Object.values(flat).reduce<TSpanMap>(
    (acc, span) => ({
      ...acc,
      [span.id || '']: Span(span),
    }),
    {}
  ),
  spans: Object.values(flat).map(rawSpan => Span(rawSpan)),
});

export default Trace;
