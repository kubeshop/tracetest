import {IRawTrace, ITrace} from '../types/Trace.types';
import Span from './Span.model';

const Trace = ({description, resourceSpans}: IRawTrace): ITrace => {
  return {
    description,
    spans: Span.createFromResourceSpanList(resourceSpans),
  };
};

export default Trace;
