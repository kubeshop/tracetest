import {NodeTypesEnum} from 'constants/DAG.constants';
import {Attributes} from 'constants/SpanAttribute.constants';
import DAGService from 'services/DAG.service';
import SpanService from 'services/Span.service';
import {INodeDataSpan, INodeDatum} from 'types/DAG.types';
import {TSpan} from 'types/Span.types';

function getNodesDatumFromSpans(spans: TSpan[]): INodeDatum<INodeDataSpan>[] {
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
      type: span.type,
      ...SpanService.getSpanNodeInfo(span),
    },
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
