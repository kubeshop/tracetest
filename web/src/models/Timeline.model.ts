import TimelineService from 'services/Timeline.service';
import {INodeDataSpan} from 'types/Timeline.types';
import Span from './Span.model';

function getNodesDataFromSpans(spans: Span[]): INodeDataSpan[] {
  return spans.map(span => ({
    duration: span.duration,
    endTime: span.endTime,
    id: span.id,
    kind: span.kind,
    name: span.name,
    parentId: span.parentId ? span.parentId : undefined,
    service: span.service,
    startTime: span.startTime,
    system: span.system,
    type: span.type,
  }));
}

function Timeline(spans: Span[]) {
  const nodesData = getNodesDataFromSpans(spans);
  return TimelineService.getNodes(nodesData);
}

export default Timeline;
