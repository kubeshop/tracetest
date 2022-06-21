import {Steps} from 'components/GuidedTour/traceStepList';
import React, {useCallback, useMemo} from 'react';

import ReactFlow, {Node} from 'react-flow-renderer';
import {useDAGChart} from 'hooks/useDAGChart';
import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import SpanService from 'services/Span.service';
// import {TSpan} from 'types/Span.types';
import TraceNode from 'components/TraceNode';
import {IDiagramComponentProps} from 'components/Diagram/Diagram';
import * as S from './DAG.styled';
import Controls from './Controls';

const {onClickSpan} = TraceDiagramAnalyticsService;

// Important to define the nodeTypes outside of the component to prevent re-renderings
const nodeTypes = {span: TraceNode};

const DAG = ({spanList = [], affectedSpans, onSelectSpan, matchedSpans, selectedSpan}: IDiagramComponentProps) => {
  const onNodeClick = useCallback(
    (event, {id}: Node) => {
      onClickSpan(id);
      if (onSelectSpan) onSelectSpan(id);
    },
    [onSelectSpan]
  );

  const nodeList = useMemo(() => SpanService.getNodeListFromSpanList(spanList), [spanList]);

  const {nodes, edges} = useDAGChart<{label: string}>(
    nodeList,
    affectedSpans,
    matchedSpans,
    selectedSpan,
    onSelectSpan
  );

  console.log('### nodes', nodes);
  console.log('### edges', edges);

  return (
    <S.Container
      data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Graph)}
      $showAffected={affectedSpans.length > 0}
      data-cy="diagram-dag"
    >
      <Controls onSelectSpan={onSelectSpan!} />
      <ReactFlow
        defaultEdges={edges}
        defaultNodes={nodes}
        deleteKeyCode={null}
        fitView
        multiSelectionKeyCode={null}
        nodesConnectable={false}
        nodeTypes={nodeTypes}
        onNodeClick={onNodeClick}
        selectionKeyCode={null}
      />
    </S.Container>
  );
};

export default DAG;
