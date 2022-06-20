import {TNode} from 'services/DAG.service';

export enum TraceNodes {
  Skeleton = 'SkeletonNode',
  TraceNode = 'TraceNode',
}

export const skeletonNodeList: TNode<{}>[] = [
  {
    id: '1',
    parentIds: [],
    data: {},
    type: TraceNodes.Skeleton,
  },
  {
    id: '2',
    parentIds: ['1'],
    data: {},
    type: TraceNodes.Skeleton,
  },
  {
    id: '3',
    parentIds: ['1'],
    data: {},
    type: TraceNodes.Skeleton,
  },
  {
    id: '4',
    parentIds: ['2'],
    data: {},
    type: TraceNodes.Skeleton,
  },
  {
    id: '5',
    parentIds: ['1'],
    data: {},
    type: TraceNodes.Skeleton,
  },
];

export const strokeColor = '#C9CEDB';
