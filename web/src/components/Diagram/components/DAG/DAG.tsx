import {Steps} from 'components/GuidedTour/traceStepList';
import React, {useCallback, useMemo} from 'react';
import ReactFlow, {Background, FlowElement} from 'react-flow-renderer';
import {TraceNodes} from 'constants/Diagram.constants';
import {useDAGChart} from 'hooks/useDAGChart';
import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {TSpan} from 'types/Span.types';
import TraceNode from 'components/TraceNode';
import {IDiagramProps} from 'components/Diagram/Diagram';
import * as S from './DAG.styled';
import Controls from './Controls';

const {onClickSpan} = TraceDiagramAnalyticsService;

const DAG: React.FC<IDiagramProps> = ({
  affectedSpans,
  trace: {spans = []},
  selectedSpan,
  onSelectSpan,
}): JSX.Element => {
  const handleElementClick = useCallback(
    (event, {id}: FlowElement) => {
      onClickSpan(id);
      if (onSelectSpan) onSelectSpan(id);
    },
    [onSelectSpan]
  );

  const nodeList = useMemo(
    () =>
      spans.map(span => ({
        id: span.id,
        parentIds: span.parentId ? [span.parentId] : [],
        data: span,
        type: TraceNodes.TraceNode,
      })),
    [spans]
  );

  const elementList = useDAGChart<TSpan>(nodeList, selectedSpan, onSelectSpan);

  return (
    <S.Container
      data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Graph)}
      $showAffected={affectedSpans.length > 0}
      data-cy="diagram-dag"
    >
      <Controls onSelectSpan={onSelectSpan!} />
      <ReactFlow
        nodeTypes={{TraceNode}}
        defaultZoom={0.5}
        elements={elementList}
        onElementClick={handleElementClick}
        onLoad={instance => setTimeout(() => instance.fitView(), 0)}
      >
        <Background gap={4} size={1} color="#FBFBFF" />
      </ReactFlow>
    </S.Container>
  );
};

export default DAG;
