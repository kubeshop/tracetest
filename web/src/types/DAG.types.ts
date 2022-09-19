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
  duration: string;
  id: string;
  isMatched: boolean;
  kind: SpanKind;
  name: string;
  programmingLanguage: string;
  service: string;
  startTime: number;
  system: string;
  totalAttributes: number;
  type: SemanticGroupNames;
}
