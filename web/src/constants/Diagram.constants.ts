import {IDAGNode} from '../hooks/useDAGChart';

export const skeletonNodeList: IDAGNode<{}>[] = [
  {
    id: '1',
    parentIds: [],
    data: {},
  },
  {
    id: '2',
    parentIds: ['1'],
    data: {},
  },
  {
    id: '3',
    parentIds: ['1'],
    data: {},
  },
  {
    id: '4',
    parentIds: ['2'],
    data: {},
  },
  {
    id: '5',
    parentIds: ['1'],
    data: {},
  },
];

export enum TraceNodes {
  Skeleton = 'SkeletonNode',
  TraceNode = 'TraceNode',
}

export const strokeColor = '#C9CEDB';
