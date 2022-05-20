import {TRawTrace, TTrace} from '../types/Trace.types';
import Span from './Span.model';

const Trace = ({traceId = '', flat = {}}: TRawTrace): TTrace => {
  return {
    traceId,
    spans: Object.values(flat).map(rawSpan => Span(rawSpan)),
  };
};

export default Trace;
