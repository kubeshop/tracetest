import {AxisScale} from '@visx/axis';
import {Group} from '@visx/group';

import {AxisOffset, BaseLeftPadding, NodeHeight, NodeOverlayHeight} from 'constants/Timeline.constants';
import useSpanData from 'hooks/useSpanData';
import {TNode} from 'types/Timeline.types';
import Collapse from './Collapse';
import Connector from './Connector';
import Label from './Label';
import * as S from './Timeline.styled';

interface IProps {
  index: number;
  indexParent: number;
  isCollapsed?: boolean;
  isMatched?: boolean;
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
  isCollapsed = false,
  isMatched = false,
  isSelected = false,
  minStartTime,
  node,
  onClick,
  onCollapse,
  xScale,
}: IProps) => {
  const {testSpecs} = useSpanData(node.data.id);
  const isParent = Boolean(node.children);
  const hasParent = indexParent !== -1;
  const positionTop = index * NodeHeight;
  const durationWidth = node.data.endTime - node.data.startTime;
  const durationX = node.data.startTime - minStartTime;
  const leftPadding = node.depth * BaseLeftPadding;
  const className = isMatched ? 'matched' : '';

  return (
    <Group className={className} id={node.data.id} left={0} top={positionTop}>
      {hasParent && <Connector distance={index - indexParent} leftPadding={leftPadding} />}

      <Group left={0} onClick={() => onClick(node.data.id)} top={0}>
        <S.RectOverlay height={NodeOverlayHeight} rx={2} x={0} y={0} $isMatched={isMatched} $isSelected={isSelected} />
      </Group>

      <Group left={leftPadding} top={8}>
        <Label
          duration={node.data.duration}
          kind={node.data.kind}
          name={node.data.name}
          service={node.data.service}
          system={node.data.system}
          totalFailedChecks={testSpecs?.failed?.length}
          totalPassedChecks={testSpecs?.passed?.length}
          type={node.data.type}
        />

        {isParent && (
          <Collapse id={node.data.id} isCollapsed={isCollapsed} onCollapse={onCollapse} totalChildren={node.children} />
        )}

        <S.RectDurationGuideline rx={3} x={50} y={46} />
      </Group>

      <S.RectDuration
        rx={3}
        width={Math.ceil(xScale(durationWidth)?.valueOf() ?? 0)}
        x={Math.ceil(xScale(durationX)?.valueOf() ?? 0) + AxisOffset}
        y={52}
        $type={node.data.type}
      />
    </Group>
  );
};

export default SpanNode;
