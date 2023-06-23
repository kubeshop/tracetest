import {NodeTypesEnum} from 'constants/DAG.constants';

export interface INodeDatum<T> {
  data: T;
  id: string;
  parentIds: string[];
  type: NodeTypesEnum;
}

export interface INodeDataSpan {
  id: string;
  isMatched: boolean;
  startTime: number;
}
