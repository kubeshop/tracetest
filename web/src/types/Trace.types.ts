import {IResourceSpan} from './Span.types';

export type ITrace = {
  resourceSpans: Array<IResourceSpan>;
  description: string;
};
