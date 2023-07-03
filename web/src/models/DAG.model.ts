import {NodeTypesEnum} from 'constants/Visualization.constants';
import DAGService from 'services/DAG.service';
import {INodeDataSpan, INodeDatum} from 'types/DAG.types';
import Span from './Span.model';

function getNodesDatumFromSpans(spans: Span[], type: NodeTypesEnum): INodeDatum<INodeDataSpan>[] {
  return spans.map(span => ({
    data: {id: span.id, isMatched: false, startTime: span.startTime},
    id: span.id,
    parentIds: span.parentId ? [span.parentId] : [],
    type,
  }));
}

function DAG(spans: Span[], type: NodeTypesEnum) {
  const nodesDatum = getNodesDatumFromSpans(spans, type).sort((a, b) => {
    if (b.data.startTime !== a.data.startTime) return b.data.startTime - a.data.startTime;
    if (b.id < a.id) return -1;
    if (b.id > a.id) return 1;
    return 0;
  });
  return DAGService.getEdgesAndNodes(nodesDatum);
}

export default DAG;
