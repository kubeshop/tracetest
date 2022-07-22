import {AxisScale} from '@visx/axis';
import {Group} from '@visx/group';

import {AxisOffset, BaseLeftPadding, NodeHeight, NodeOverlayHeight} from 'constants/Timeline.constants';
import {TNode} from 'types/Timeline.types';
import Collapse from './Collapse';
import Connector from './Connector';
import Label from './Label';
import * as S from './Timeline.styled';

interface IProps {
  index: number;
  indexParent: number;
  isAffected?: boolean;
  isCollapsed?: boolean;
  isSelected?: boolean;
  minStartTime: number;
  node: TNode;
  onClick(id: string): void;
  onCollapse(id: string): void;
  xScale: AxisScale;
}

const SpanNode = ({
  index,
  indexParent,
  isAffected = false,
  isCollapsed = false,
  isSelected = false,
  minStartTime,
  node,
  onClick,
  onCollapse,
  xScale,
}: IProps) => {
  const isParent = Boolean(node.children);
  const hasParent = indexParent !== -1;
  const positionTop = index * NodeHeight;
  const durationWidth = node.data.endTime - node.data.startTime;
  const durationX = node.data.startTime - minStartTime;
  const leftPadding = node.depth * BaseLeftPadding;

  return (
    <Group left={0} top={positionTop}>
      {hasParent && <Connector distance={index - indexParent} leftPadding={leftPadding} />}

      <Group left={0} onClick={() => onClick(node.data.id)} top={0}>
        <S.RectOverlay
          height={NodeOverlayHeight}
          rx={2}
          x={0}
          y={0}
          $isAffected={isAffected}
          $isSelected={isSelected}
        />
      </Group>

      <Group left={leftPadding} top={8}>
        <Label
          duration={node.data.duration}
          kind={node.data.kind}
          name={node.data.name}
          service={node.data.service}
          system={node.data.system}
          type={node.data.type}
        />

        {isParent && (
          <Collapse id={node.data.id} isCollapsed={isCollapsed} onCollapse={onCollapse} totalChildren={node.children} />
        )}

        <S.RectDurationGuideline rx={3} x={50} y={46} />
      </Group>

      <S.RectDuration
        rx={3}
        width={xScale(durationWidth) as number}
        x={(xScale(durationX) as number) + AxisOffset}
        y={52}
        $type={node.data.type}
      />
    </Group>
  );
};

export default SpanNode;
