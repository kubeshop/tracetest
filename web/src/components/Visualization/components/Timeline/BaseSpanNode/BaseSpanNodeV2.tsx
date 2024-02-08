import Span from 'models/Span.model';
import Connector from './ConnectorV2';
import {IPropsComponent} from '../SpanNodeFactoryV2';
import {useTimeline} from '../Timeline.provider';
import * as S from '../TimelineV2.styled';

const BaseLeftPadding = 16; // TODO: Move to Timeline.constants

function toPercent(value: number) {
  return `${(value * 100).toFixed(1)}%`;
}

function getHintSide(viewStart: number, viewEnd: number) {
  return viewStart > 1 - viewEnd ? 'left' : 'right';
}

interface IProps extends IPropsComponent {
  header?: React.ReactNode;
  span: Span;
}

const BaseSpanNode = ({header, index, node, span, style}: IProps) => {
  const {getScale, matchedSpans, onNodeClick, selectedSpan} = useTimeline();
  const {start: viewStart, end: viewEnd} = getScale(span.startTime, span.endTime);
  const hintSide = getHintSide(viewStart, viewEnd);
  const isSelected = selectedSpan === node.data.id;
  const isMatched = matchedSpans.includes(node.data.id);
  const leftPadding = node.depth * BaseLeftPadding;

  return (
    <div style={style}>
      <S.Row
        onClick={() => onNodeClick(node.data.id)}
        $isEven={index % 2 === 0}
        $isMatched={isMatched}
        $isSelected={isSelected}
      >
        <S.Col>
          <S.Header>
            <Connector hasParent={!!node.data.parentId} leftPadding={leftPadding} totalChildren={node.children} />
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
