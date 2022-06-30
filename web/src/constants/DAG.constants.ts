import {INodeDatum} from 'types/DAG.types';

export enum NodeTypesEnum {
  Span = 'span',
  Skeleton = 'skeleton',
}

export enum Colors {
  Default = '#CDD1DB',
  Selected = '#2d62ff',
}

export const skeletonNodesDatum: INodeDatum<{}>[] = [
  {
    id: '1',
    parentIds: [],
    data: {},
    type: NodeTypesEnum.Skeleton,
  },
  {
    id: '2',
    parentIds: ['1'],
    data: {},
    type: NodeTypesEnum.Skeleton,
  },
  {
    id: '3',
    parentIds: ['1'],
    data: {},
    type: NodeTypesEnum.Skeleton,
  },
  {
    id: '4',
    parentIds: ['2'],
    data: {},
    type: NodeTypesEnum.Skeleton,
  },
  {
    id: '5',
    parentIds: ['1'],
    data: {},
    type: NodeTypesEnum.Skeleton,
  },
];
