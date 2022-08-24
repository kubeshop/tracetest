import ReactFlow, {MiniMap} from 'react-flow-renderer';

import {IDiagramComponentProps} from 'components/Diagram/Diagram';
import {Steps} from 'components/GuidedTour/traceStepList';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {useDAG} from 'providers/DAG';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import Controls from './Controls';
import * as S from './DAG.styled';
import SpanNode from './SpanNode';

/** Important to define the nodeTypes outside the component to prevent re-renderings */
const nodeTypes = {span: SpanNode};

const DAG = ({affectedSpans}: IDiagramComponentProps) => {
  const {edges, isMiniMapActive, nodes, onMiniMapToggle, onNodesChange, onNodeClick} = useDAG();
  const {isOpen} = useTestSpecForm();

  return (
    <S.Container
      data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Graph)}
      $showAffected={affectedSpans.length > 0 || isOpen}
      data-cy="diagram-dag"
    >
      <Controls isMiniMapActive={isMiniMapActive} onMiniMapToggle={onMiniMapToggle} />
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
  );
};

export default DAG;
