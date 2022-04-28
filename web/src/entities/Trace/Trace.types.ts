import {TResourceSpan} from '../Span/Span.types';

export type TTrace = {
  resourceSpans: Array<TResourceSpan>;
  description: string;
};
