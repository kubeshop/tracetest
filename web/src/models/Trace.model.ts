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

const Trace = ({traceId = '', flat: rawFlat = {}, tree = {}} = defaultTrace): Trace => {
  const flat: TSpanMap = {};
  const spans = Object.values(rawFlat).map(raw => {
    const span = Span(raw);
    flat[span.id || ''] = span;

    return span;
  });

  return {
    traceId,
    rootSpan: Span(tree),
    flat,
    spans,
  };
};

export default Trace;
