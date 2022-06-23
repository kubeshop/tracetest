import {NodeTypesEnum} from 'constants/DAG.constants';
import DAGService from 'services/DAG.service';
import SpanService from 'services/Span.service';
import {INodeDataSpan, INodeDatum} from 'types/DAG.types';
import {TSpan} from 'types/Span.types';

function getNodesDatumFromSpans(spans: TSpan[]): INodeDatum<INodeDataSpan>[] {
  return spans.map(span => ({
    data: {name: span.name, type: span.type, isAffected: false, isMatched: false, ...SpanService.getSpanNodeInfo(span)},
    id: span.id,
    parentIds: span.parentId ? [span.parentId] : [],
    type: NodeTypesEnum.Span,
  }));
}

function DAG(spans: TSpan[]) {
  const nodesDatum = getNodesDatumFromSpans(spans);
  return DAGService.getEdgesAndNodes(nodesDatum);
}

export default DAG;
