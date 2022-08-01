import React, {useEffect, useMemo} from 'react';
import {Typography} from 'antd';
import ReactFlow from 'react-flow-renderer';

import {TRACE_DOCUMENTATION_URL} from 'constants/Common.constants';
import {skeletonNodesDatum} from 'constants/DAG.constants';
import DAGService from 'services/DAG.service';
import * as S from './SkeletonDiagram.styled';
import SkeletonNode from './SkeletonNode';

export interface IProps {
  onClearAffectedSpans(): void;
  onClearSelectedSpan(): void;
}

/** Important to define the nodeTypes outside of the component to prevent re-renderings */
const nodeTypes = {skeleton: SkeletonNode};

const SkeletonDiagram = ({onClearAffectedSpans, onClearSelectedSpan}: IProps) => {
  const {edges, nodes} = useMemo(() => DAGService.getEdgesAndNodes(skeletonNodesDatum), []);

  useEffect(() => {
    onClearAffectedSpans();
    onClearSelectedSpan();
  }, []);

  return (
    <S.Container data-cy="skeleton-diagram">
      <S.SkeletonDiagramMessage>
        <Typography.Title level={3} type="secondary">
          We are working on your trace…
        </Typography.Title>
        <Typography.Text type="secondary">
          Want to know more about traces? head to the official{' '}
          <a href={TRACE_DOCUMENTATION_URL} target="_blank">
            Open Telemetry Documentation
          </a>
        </Typography.Text>
      </S.SkeletonDiagramMessage>

      <ReactFlow
        defaultEdges={edges}
        defaultNodes={nodes}
        deleteKeyCode={null}
        fitView
        multiSelectionKeyCode={null}
        nodesConnectable={false}
        nodeTypes={nodeTypes}
        selectionKeyCode={null}
      />
    </S.Container>
  );
};

export default SkeletonDiagram;
