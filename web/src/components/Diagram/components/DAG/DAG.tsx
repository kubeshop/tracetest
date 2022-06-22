import React from 'react';
import ReactFlow from 'react-flow-renderer';

import {IDiagramComponentProps} from 'components/Diagram/Diagram';
import {Steps} from 'components/GuidedTour/traceStepList';
import TraceNode from 'components/TraceNode';
import {useDAG} from 'providers/DAG';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import Controls from './Controls';
import * as S from './DAG.styled';

/** Important to define the nodeTypes outside of the component to prevent re-renderings */
const nodeTypes = {span: TraceNode};

const DAG = ({affectedSpans, onSelectSpan}: IDiagramComponentProps) => {
  const {edges, nodes, onNodesChange, onNodeClick} = useDAG();

  return (
    <S.Container
      data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Graph)}
      $showAffected={affectedSpans.length > 0}
      data-cy="diagram-dag"
    >
      <Controls onSelectSpan={onSelectSpan!} />
      <ReactFlow
        edges={edges}
        nodes={nodes}
        deleteKeyCode={null}
        fitView
        multiSelectionKeyCode={null}
        nodesConnectable={false}
        nodeTypes={nodeTypes}
        onNodeClick={onNodeClick}
        onNodesChange={onNodesChange}
        selectionKeyCode={null}
      />
    </S.Container>
  );
};

export default DAG;
