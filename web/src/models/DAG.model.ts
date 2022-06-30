import {NodeTypesEnum} from 'constants/DAG.constants';
import {SemanticGroupNamesToSystem} from 'constants/SemanticGroupNames.constants';
import {Attributes} from 'constants/SpanAttribute.constants';
import DAGService from 'services/DAG.service';
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
      system: span.attributes?.[SemanticGroupNamesToSystem[span.type]]?.value ?? '',
      totalAttributes: span.attributeList.length,
      type: span.type,
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
