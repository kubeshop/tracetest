import Span from 'models/Span.model';
import Connector from './Connector';
import {IPropsComponent} from '../SpanNodeFactory';
import {useTimeline} from '../Timeline.provider';
import * as S from '../Timeline.styled';

function toPercent(value: number) {
  return `${(value * 100).toFixed(1)}%`;
}

function getHintSide(viewStart: number, viewEnd: number) {
  return viewStart > 1 - viewEnd ? 'left' : 'right';
}

interface IProps extends IPropsComponent {
  span: Span;
}

const BaseSpanNode = ({index, node, span, style}: IProps) => {
  const {collapsedSpans, getScale, matchedSpans, onSpanCollapse, onSpanClick, selectedSpan} = useTimeline();
  const {start: viewStart, end: viewEnd} = getScale(span.startTime, span.endTime);
  const hintSide = getHintSide(viewStart, viewEnd);
  const isSelected = selectedSpan === node.data.id;
  const isMatched = matchedSpans.includes(node.data.id);
  const isCollapsed = collapsedSpans.includes(node.data.id);

  return (
    <div style={style}>
      <S.Row
        onClick={() => onSpanClick(node.data.id)}
        $isEven={index % 2 === 0}
        $isMatched={isMatched}
        $isSelected={isSelected}
      >
        <S.Col>
          <S.Header>
            <Connector
              hasParent={!!node.data.parentId}
              id={node.data.id}
              isCollapsed={isCollapsed}
              nodeDepth={node.depth}
              onCollapse={onSpanCollapse}
              totalChildren={node.children}
            />
            <S.NameContainer>
              <S.Title>{span.name}</S.Title>
            </S.NameContainer>
          </S.Header>
          <S.Separator />
        </S.Col>

        <S.ColDuration>
          <S.SpanBar $type={span.type} style={{left: toPercent(viewStart), width: toPercent(viewEnd - viewStart)}}>
            <S.SpanBarLabel $side={hintSide}>{span.duration}</S.SpanBarLabel>
          </S.SpanBar>
        </S.ColDuration>
      </S.Row>
    </div>
  );
};

export default BaseSpanNode;
