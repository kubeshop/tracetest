import TimelineService from 'services/Timeline.service';
import {TSpan} from 'types/Span.types';
import {INodeDataSpan} from 'types/Timeline.types';

function getNodesDataFromSpans(spans: TSpan[]): INodeDataSpan[] {
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

function Timeline(spans: TSpan[]) {
  const nodesData = getNodesDataFromSpans(spans);
  return TimelineService.getNodes(nodesData);
}

export default Timeline;
