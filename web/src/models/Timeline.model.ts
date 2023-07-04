import {NodeTypesEnum} from 'constants/Visualization.constants';
import TimelineService from 'services/Timeline.service';
import {INodeDataSpan} from 'types/Timeline.types';
import Span from './Span.model';

function getNodesDataFromSpans(spans: Span[]): INodeDataSpan[] {
  return spans.map(span => ({
    endTime: span.endTime,
    id: span.id,
    parentId: span.parentId ? span.parentId : undefined,
    startTime: span.startTime,
  }));
}

function Timeline(spans: Span[], type: NodeTypesEnum) {
  const nodesData = getNodesDataFromSpans(spans);
  return TimelineService.getNodes(nodesData, type);
}

export default Timeline;
