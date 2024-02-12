import {NodeTypesEnum} from 'constants/Visualization.constants';
import Span from 'models/Span.model';
import {useRef} from 'react';
import {FixedSizeList as List} from 'react-window';
import NavigationWrapper from './NavigationWrapper';
import TimelineProvider from './Timeline.provider';
import ListWrapper from './ListWrapper';

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
      <NavigationWrapper />
      <ListWrapper listRef={listRef} />
    </TimelineProvider>
  );
};

export default Timeline;
