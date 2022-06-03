import React, {useEffect, useMemo} from 'react';
import {Typography} from 'antd';
import ReactFlow, {ArrowHeadType, Background, Elements, Position} from 'react-flow-renderer';
import SkeletonNode from './SkeletonNode';
import * as S from './SkeletonDiagram.styled';
import {useDAGChart} from '../../hooks/useDAGChart';
import {TRACE_DOCUMENTATION_URL} from '../../constants/Common.constants';
import {skeletonNodeList, strokeColor, TraceNodes} from '../../constants/Diagram.constants';

export type SkeletonElementList = Elements<{}>;

export interface IProps {
  onSelectSpan?(spanId: string): void;
}

const SkeletonDiagram = ({onSelectSpan}: IProps) => {
  const {dag} = useDAGChart(skeletonNodeList);

  useEffect(() => {
    if (onSelectSpan) onSelectSpan('');
  }, []);

  const dagElementList = useMemo<SkeletonElementList>(() => {
    const dagNodeList: SkeletonElementList =
      dag?.descendants().map(({data, x, y}) => ({
        id: data.id,
        type: TraceNodes.Skeleton,
        position: {x: x!, y: parseFloat(String(y))},
        data,
        sourcePosition: Position.Top,
      })) || [];

    dag?.links().forEach(({source, target}) => {
      dagNodeList.push({
        id: `${source.data.id}_${target.data.id}`,
        source: source.data.id,
        target: target.data.id,
        animated: true,
        arrowHeadType: ArrowHeadType.ArrowClosed,
        style: {stroke: strokeColor},
      });
    });

    return dagNodeList;
  }, [dag]);

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
        elements={dagElementList}
        onLoad={instance => setTimeout(() => instance.fitView(), 0)}
      >
        <Background gap={4} size={1} color="#FBFBFF" />
      </ReactFlow>
    </S.Container>
  );
};

export default SkeletonDiagram;
