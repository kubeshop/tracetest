import {Steps} from 'components/GuidedTour/traceStepList';
import React, {useCallback, useMemo} from 'react';
import ReactFlow, {Background, FlowElement} from 'react-flow-renderer';
import {useDAGChart} from 'hooks/useDAGChart';
import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import SpanService from 'services/Span.service';
import {TSpan} from 'types/Span.types';
import TraceNode from 'components/TraceNode';
import {IDiagramComponentProps} from 'components/Diagram/Diagram';
import * as S from './DAG.styled';
import Controls from './Controls';

const {onClickSpan} = TraceDiagramAnalyticsService;

const DAG = ({spanList = [], affectedSpans, onSelectSpan, matchedSpans, selectedSpan}: IDiagramComponentProps) => {
  const handleElementClick = useCallback(
    (event, {id}: FlowElement) => {
      onClickSpan(id);
      if (onSelectSpan) onSelectSpan(id);
    },
    [onSelectSpan]
  );

  const nodeList = useMemo(
    () => SpanService.getNodeListFromSpanList(spanList, affectedSpans, matchedSpans),
    [affectedSpans, matchedSpans, spanList]
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
