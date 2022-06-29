import {NodeTypesEnum} from 'constants/DAG.constants';
import {Attributes} from 'constants/SpanAttribute.constants';
import DAGService from 'services/DAG.service';
import SpanService from 'services/Span.service';
import {INodeDataSpan, INodeDatum} from 'types/DAG.types';
import {TSpan, TSpansResult} from 'types/Span.types';

function getNodesDatumFromSpans(spans: TSpan[], spansResult: TSpansResult): INodeDatum<INodeDataSpan>[] {
  return spans.map(span => ({
    data: {
      duration: span.duration,
      id: span.id,
      isAffected: false,
      isMatched: false,
      kind: span.kind,
      name: span.name,
      programmingLanguage: span.attributes?.[Attributes.TELEMETRY_SDK_LANGUAGE]?.value ?? '',
      serviceName: span.attributes?.[Attributes.SERVICE_NAME]?.value ?? '',
      totalAttributes: span.attributeList.length,
      totalChecksFailed: spansResult?.[span.id]?.failed ?? 0,
      totalChecksPassed: spansResult?.[span.id]?.passed ?? 0,
      type: span.type,
      ...SpanService.getSpanNodeInfo(span),
    },
    id: span.id,
    parentIds: span.parentId ? [span.parentId] : [],
    type: NodeTypesEnum.Span,
  }));
}

function DAG(spans: TSpan[], spansResult: TSpansResult) {
  const nodesDatum = getNodesDatumFromSpans(spans, spansResult);
  return DAGService.getEdgesAndNodes(nodesDatum);
}

export default DAG;
