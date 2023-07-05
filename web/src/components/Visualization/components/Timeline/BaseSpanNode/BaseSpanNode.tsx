import {Group} from '@visx/group';

import {AxisOffset, BaseLeftPadding, NodeHeight, NodeOverlayHeight} from 'constants/Timeline.constants';
import Span from 'models/Span.model';
import Collapse from './Collapse';
import Connector from './Connector';
import Label from './Label';
import {IPropsComponent} from '../SpanNodeFactory';
import * as S from '../Timeline.styled';

interface IProps extends IPropsComponent {
  className?: string;
  header?: React.ReactNode;
  span: Span;
}

const BaseSpanNode = ({
  className,
  header,
  index,
  indexParent,
  isCollapsed = false,
  isMatched = false,
  isSelected = false,
  minStartTime,
  node,
  onClick,
  onCollapse,
  span,
  xScale,
}: IProps) => {
  const isParent = Boolean(node.children);
  const hasParent = indexParent !== -1;
  const positionTop = index * NodeHeight;
  const durationWidth = span.endTime - span.startTime;
  const durationX = span.startTime - minStartTime;
  const leftPadding = node.depth * BaseLeftPadding;

  return (
    <Group className={className} id={node.data.id} left={0} top={positionTop}>
      {hasParent && <Connector distance={index - indexParent} leftPadding={leftPadding} />}

      <Group left={0} onClick={() => onClick(node.data.id)} top={0}>
        <S.RectOverlay height={NodeOverlayHeight} rx={2} x={0} y={0} $isMatched={isMatched} $isSelected={isSelected} />
      </Group>

      <Group left={leftPadding} top={8}>
        <Label
          duration={span.duration}
          header={header}
          kind={span.kind}
          name={span.name}
          service={span.service}
          system={span.system}
          type={span.type}
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
        $type={span.type}
      />
    </Group>
  );
};

export default BaseSpanNode;
