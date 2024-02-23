import {NodeTypesEnum} from 'constants/Visualization.constants';
import Span from 'models/Span.model';
import {useRef} from 'react';
import {FixedSizeList as List} from 'react-window';
import ListWrapper from './ListWrapper';
import NavigationWrapper from './NavigationWrapper';
import TimelineProvider from './Timeline.provider';
import VerticalResizerProvider from './VerticalResizer.provider';

export interface IProps {
  nodeType: NodeTypesEnum;
  spans: Span[];
  onNavigate(spanId: string): void;
  onClick(spanId: string): void;
  matchedSpans: string[];
  selectedSpan: string;
}

const Timeline = ({nodeType, spans, onClick, onNavigate, matchedSpans, selectedSpan}: IProps) => {
  const listRef = useRef<List>(null);

  return (
    <TimelineProvider
      onClick={onClick}
      onNavigate={onNavigate}
      matchedSpans={matchedSpans}
      selectedSpan={selectedSpan}
      listRef={listRef}
      nodeType={nodeType}
      spans={spans}
    >
      <VerticalResizerProvider>
        <NavigationWrapper />
        <ListWrapper listRef={listRef} />
      </VerticalResizerProvider>
    </TimelineProvider>
  );
};

export default Timeline;
