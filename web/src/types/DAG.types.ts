import {NodeTypesEnum} from 'constants/DAG.constants';
import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {SpanKind} from 'constants/Span.constants';

export interface INodeDatum<T> {
  data: T;
  id: string;
  parentIds: string[];
  type: NodeTypesEnum;
}

export interface INodeDataSpan {
  duration: number;
  heading: string;
  id: string;
  isAffected: boolean;
  isMatched: boolean;
  kind: SpanKind;
  name: string;
  primary: string;
  programmingLanguage: string;
  serviceName: string;
  totalAttributes: number;
  totalChecksFailed: number;
  totalChecksPassed: number;
  type: SemanticGroupNames;
}
