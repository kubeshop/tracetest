import {IResourceSpan, ISpan} from './Span.types';

export type IRawTrace = {
  resourceSpans?: Array<IResourceSpan>;
  description?: string;
};

export type ITrace = {
  spans: ISpan[];
  description: string;
};
