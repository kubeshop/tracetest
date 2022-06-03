import {ZoomInOutlined, ZoomOutOutlined} from '@ant-design/icons';

import {Steps} from 'components/GuidedTour/traceStepList';
import React, {useCallback, useEffect, useMemo} from 'react';
import ReactFlow, {
  ArrowHeadType,
  Background,
  Elements,
  FlowElement,
  Position,
  useZoomPanHelper,
} from 'react-flow-renderer';
import {strokeColor, TraceNodes} from '../../../constants/Diagram.constants';
import {IDAGNode, useDAGChart} from '../../../hooks/useDAGChart';
import TraceDiagramAnalyticsService from '../../../services/Analytics/TraceDiagramAnalytics.service';
import GuidedTourService, {GuidedTours} from '../../../services/GuidedTour.service';
import {TSpan} from '../../../types/Span.types';
import TraceNode from '../../TraceNode';
import {IDiagramProps} from '../Diagram';
import * as S from './DAG.styled';

type TElementList = Elements<TSpan>;

const {onClickSpan} = TraceDiagramAnalyticsService;

const Diagram: React.FC<IDiagramProps> = ({
  affectedSpans,
  trace: {spans = []},
  selectedSpan,
  onSelectSpan,
}): JSX.Element => {
  const {zoomIn, zoomOut} = useZoomPanHelper();

  const nodeList = useMemo<IDAGNode<TSpan>[]>(
    () =>
      spans.map(span => ({
        id: span.id,
        parentIds: span.parentId ? [span.parentId] : [],
        data: span,
      })),
    [spans]
  );

  const {dag} = useDAGChart(nodeList);

  const handleElementClick = useCallback(
    (event, {id}: FlowElement) => {
      onClickSpan(id);
      if (onSelectSpan) onSelectSpan(id);
    },
    [onSelectSpan]
  );

  useEffect(() => {
    if (dag) {
      const [dagNode] = dag.descendants();
      const node = nodeList.find(({id}) => id === dagNode?.data.id);

      if (!selectedSpan && node && onSelectSpan) onSelectSpan(node.id);
    }
  }, [dag, nodeList, onSelectSpan, selectedSpan]);

  const dagElements = useMemo<TElementList>(() => {
    if (dag) {
      const dagNodeList: TElementList = dag.descendants().map(({data: {data}, x, y}) => {
        return {
          id: data.id,
          type: TraceNodes.TraceNode,
          data,
          position: {x: x!, y: parseFloat(String(y))},
          selected: data.id === selectedSpan?.id,
          sourcePosition: Position.Top,
          className: affectedSpans.includes(data.id) ? 'affected' : '',
        };
      });

      dag.links().forEach(({source, target}) => {
        dagNodeList.push({
          id: `${source.data.id}_${target.data.id}`,
          source: source.data.id,
          target: target.data.id,
          data: source.data.data,
          labelShowBg: false,
          animated: false,
          arrowHeadType: ArrowHeadType.ArrowClosed,
          style: {stroke: strokeColor},
        });
      });

      return dagNodeList;
    }

    return [];
  }, [dag, selectedSpan?.id, affectedSpans]);

  return (
    <S.Container
      data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Graph)}
      $showAffected={affectedSpans.length > 0}
      data-cy="diagram-dag"
    >
      <S.Controls>
        <S.ZoomButton icon={<ZoomInOutlined />} onClick={() => zoomIn()} type="text" />
        <S.ZoomButton icon={<ZoomOutOutlined />} onClick={() => zoomOut()} type="text" />
      </S.Controls>
      <ReactFlow
        nodeTypes={{TraceNode}}
        defaultZoom={0.5}
        elements={dagElements}
        onElementClick={handleElementClick}
        onLoad={instance => setTimeout(() => instance.fitView(), 0)}
      >
        <Background gap={4} size={1} color="#FBFBFF" />
      </ReactFlow>
    </S.Container>
  );
};

export default Diagram;
