import {MouseEvent, useState} from 'react';
import ReactFlow, {Edge, MiniMap, Node, NodeChange, ReactFlowProvider} from 'react-flow-renderer';

import Actions from './Actions';
import * as S from './DAG.styled';
import TestSpanNode from './TestSpanNode/TestSpanNode';
import TraceSpanNode from './TraceSpanNode/TraceSpanNode';

/** Important to define the nodeTypes outside the component to prevent re-renderings */
const nodeTypes = {traceSpan: TraceSpanNode, testSpan: TestSpanNode};

interface IProps {
  edges: Edge[];
  isMatchedMode: boolean;
  matchedSpans: string[];
  nodes: Node[];
  onNavigateToSpan(spanId: string): void;
  onNodesChange(changes: NodeChange[]): void;
  onNodeClick(event: MouseEvent, node: Node): void;
  selectedSpan: string;
}

const DAG = ({
  edges,
  isMatchedMode,
  matchedSpans,
  nodes,
  onNavigateToSpan,
  onNodesChange,
  onNodeClick,
  selectedSpan,
}: IProps) => {
  const [isMiniMapActive, setIsMiniMapActive] = useState(false);

  return (
    <ReactFlowProvider>
      <S.Container $showMatched={isMatchedMode} data-cy="diagram-dag">
        <Actions
          isMiniMapActive={isMiniMapActive}
          matchedSpans={matchedSpans}
          onMiniMapToggle={() => setIsMiniMapActive(isActive => !isActive)}
          onNavigateToSpan={onNavigateToSpan}
          selectedSpan={selectedSpan}
        />
        <ReactFlow
          edges={edges}
          nodes={nodes}
          deleteKeyCode={null}
          fitView
          minZoom={0.1}
          multiSelectionKeyCode={null}
          nodesConnectable={false}
          nodeTypes={nodeTypes}
          onNodeClick={onNodeClick}
          onNodeDragStop={onNodeClick}
          onNodesChange={onNodesChange}
          selectionKeyCode={null}
        >
          {isMiniMapActive && <MiniMap />}
        </ReactFlow>
      </S.Container>
    </ReactFlowProvider>
  );
};

export default DAG;
