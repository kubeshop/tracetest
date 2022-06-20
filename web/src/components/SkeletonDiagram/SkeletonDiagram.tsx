import React, {useEffect} from 'react';
import {Typography} from 'antd';
import ReactFlow, {Background, Elements} from 'react-flow-renderer';
import {useDAGChart} from 'hooks/useDAGChart';
import {TRACE_DOCUMENTATION_URL} from 'constants/Common.constants';
import {skeletonNodeList} from 'constants/Diagram.constants';
import {TSpan} from 'types/Span.types';
import SkeletonNode from './SkeletonNode';
import * as S from './SkeletonDiagram.styled';

export type SkeletonElementList = Elements<{}>;

export interface IProps {
  onSelectSpan?(spanId: string): void;
  selectedSpan?: TSpan;
}

const SkeletonDiagram = ({onSelectSpan, selectedSpan}: IProps) => {
  const elementList = useDAGChart(skeletonNodeList, selectedSpan);

  useEffect(() => {
    if (onSelectSpan) onSelectSpan('');
  }, []);

  return (
    <S.Container data-cy="skeleton-diagram">
      <S.SkeletonDiagramMessage>
        <Typography.Title level={5} type="secondary">
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
        nodeTypes={{SkeletonNode}}
        defaultZoom={0.5}
        elements={elementList}
        onLoad={instance => setTimeout(() => instance.fitView(), 0)}
      >
        <Background gap={4} size={1} color="#FBFBFF" />
      </ReactFlow>
    </S.Container>
  );
};

export default SkeletonDiagram;
