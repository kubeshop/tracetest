import {NodeTypesEnum} from 'constants/Visualization.constants';

export interface INode<T> {
  children: number;
  data: T;
  depth: number;
  type: NodeTypesEnum;
}

export interface INodeDataSpan {
  id: string;
  parentId: string | undefined;
  startTime: number;
  endTime: number;
}

export type TNode = INode<INodeDataSpan>;
