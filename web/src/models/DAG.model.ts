import {NodeTypesEnum} from 'constants/DAG.constants';
import {Attributes} from 'constants/SpanAttribute.constants';
import DAGService from 'services/DAG.service';
import {INodeDataSpan, INodeDatum} from 'types/DAG.types';
import Span from './Span.model';

function getNodesDatumFromSpans(spans: Span[]): INodeDatum<INodeDataSpan>[] {
  return spans.map(span => ({
    data: {
      duration: span.duration,
      id: span.id,
      isMatched: false,
      kind: span.kind,
      name: span.name,
      programmingLanguage: span.attributes?.[Attributes.TELEMETRY_SDK_LANGUAGE]?.value ?? '',
      service: span.service,
      startTime: span.startTime,
      system: span.system,
      totalAttributes: span.attributeList.length,
      type: span.type,
    },
    id: span.id,
    parentIds: span.parentId ? [span.parentId] : [],
    type: NodeTypesEnum.Span,
  }));
}

function DAG(spans: Span[]) {
  const nodesDatum = getNodesDatumFromSpans(spans).sort((a, b) => {
    if (b.data.startTime !== a.data.startTime) return b.data.startTime - a.data.startTime;
    if (b.id < a.id) return -1;
    if (b.id > a.id) return 1;
    return 0;
  });
  return DAGService.getEdgesAndNodes(nodesDatum);
}

export default DAG;
