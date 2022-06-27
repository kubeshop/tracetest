import {NodeTypesEnum} from 'constants/DAG.constants';
import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';

export interface INodeDatum<T> {
  data: T;
  id: string;
  parentIds: string[];
  type: NodeTypesEnum;
}

export interface INodeDataSpan {
  heading: string;
  name: string;
  primary: string;
  type: SemanticGroupNames;
  isAffected: boolean;
  isMatched: boolean;
}
