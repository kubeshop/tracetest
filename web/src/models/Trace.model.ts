import {TRawTrace, TTrace} from '../types/Trace.types';
import Span from './Span.model';

const Trace = ({resourceSpans = []}: TRawTrace): TTrace => {
  return {
    description: '',
    spans: Span.createFromResourceSpanList(resourceSpans),
  };
};

export default Trace;
