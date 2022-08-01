import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {SpanKind} from 'constants/Span.constants';

export interface INode<T> {
  children: number;
  data: T;
  depth: number;
}

export interface INodeDataSpan {
  duration: string;
  id: string;
  kind: SpanKind;
  name: string;
  service: string;
  system: string;
  parentId: string | undefined;
  type: SemanticGroupNames;
  startTime: number;
  endTime: number;
}

export type TNode = INode<INodeDataSpan>;
