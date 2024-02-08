import {NodeTypesEnum} from 'constants/Visualization.constants';
import Span from 'models/Span.model';
import TimelineModel from 'models/Timeline.model';
import {useMemo, useRef} from 'react';
import {FixedSizeList as List} from 'react-window';
import NavigationWrapper from './NavigationWrapper';
import SpanNodeFactory from './SpanNodeFactoryV2';
import TimelineProvider from './Timeline.provider';
import * as S from './TimelineV2.styled';

const HEADER_HEIGHT = 210;

export interface IProps {
  nodeType: NodeTypesEnum;
  spans: Span[];
}

const Timeline = ({nodeType, spans}: IProps) => {
  const listRef = useRef<List>(null);
  const nodes = useMemo(() => TimelineModel(spans, nodeType), [spans, nodeType]);

  return (
    <TimelineProvider listRef={listRef} nodes={nodes}>
      <NavigationWrapper />
      <S.Container>
        <List
          height={window.innerHeight - HEADER_HEIGHT}
          itemCount={nodes.length}
          itemData={nodes}
          itemSize={32}
          ref={listRef}
          width="100%"
        >
          {SpanNodeFactory}
        </List>
      </S.Container>
    </TimelineProvider>
  );
};

export default Timeline;
