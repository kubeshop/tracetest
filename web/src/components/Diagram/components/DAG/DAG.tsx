import {Steps} from 'components/GuidedTour/traceStepList';
import React, {useCallback, useMemo} from 'react';

import ReactFlow, {Node} from 'react-flow-renderer';
import {useDAGChart} from 'hooks/useDAGChart';
// import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
// import {TSpan} from 'types/Span.types';
import SpanService from 'services/Span.service';
import TraceNode from 'components/TraceNode';
import {IDiagramComponentProps} from 'components/Diagram/Diagram';
import * as S from './DAG.styled';
import Controls from './Controls';

// const {onClickSpan} = TraceDiagramAnalyticsService;

// Important to define the nodeTypes outside of the component to prevent re-renderings
const nodeTypes = {span: TraceNode};

const DAG = ({spanList = [], affectedSpans, onSelectSpan, matchedSpans, selectedSpan}: IDiagramComponentProps) => {
  const onNodeClick = useCallback(
    (event, {id}: Node) => {
      // onClickSpan(id);
      if (onSelectSpan) onSelectSpan(id);
    },
    [onSelectSpan]
  );

  const items = useMemo(
    () => SpanService.getNodeListFromSpanList(spanList, affectedSpans, matchedSpans),
    [affectedSpans, matchedSpans, spanList]
  );

  const {nodes, edges} = useDAGChart<{label: string}>(items, selectedSpan, onSelectSpan);

  console.log('### items', items);
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
        defaultNodes={nodes}
        defaultEdges={edges}
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
