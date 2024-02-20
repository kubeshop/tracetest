import {NodeTypesEnum} from 'constants/Visualization.constants';
import noop from 'lodash/noop';
import without from 'lodash/without';
import Span from 'models/Span.model';
import TimelineModel from 'models/Timeline.model';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {FixedSizeList as List} from 'react-window';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import TimelineService, {TScaleFunction} from 'services/Timeline.service';
import {TNode} from 'types/Timeline.types';

interface IContext {
  collapsedSpans: string[];
  getScale: TScaleFunction;
  matchedSpans: string[];
  nameColumnWidth: number;
  onNameColumnWidthChange(width: number): void;
  onSpanClick(spanId: string): void;
  onSpanCollapse(spanId: string): void;
  onSpanNavigation(spanId: string): void;
  selectedSpan: string;
  spans: TNode[];
  viewEnd: number;
  viewStart: number;
}

export const Context = createContext<IContext>({
  collapsedSpans: [],
  getScale: () => ({start: 0, end: 0}),
  matchedSpans: [],
  nameColumnWidth: 0,
  onNameColumnWidthChange: noop,
  onSpanClick: noop,
  onSpanCollapse: noop,
  onSpanNavigation: noop,
  selectedSpan: '',
  spans: [],
  viewEnd: 0,
  viewStart: 0,
});

interface IProps {
  children: React.ReactNode;
  listRef: React.RefObject<List>;
  nodeType: NodeTypesEnum;
  spans: Span[];
  onNavigate(spanId: string): void;
  onClick(spanId: string): void;
  matchedSpans: string[];
  selectedSpan: string;
}

export const useTimeline = () => useContext(Context);

const TimelineProvider = ({
  children,
  listRef,
  nodeType,
  spans,
  onClick,
  onNavigate,
  matchedSpans,
  selectedSpan,
}: IProps) => {
  const [collapsedSpans, setCollapsedSpans] = useState<string[]>([]);
  const [nameColumnWidth, setNameColumnWidth] = useState(0.15);

  const nodes = useMemo(() => TimelineModel(spans, nodeType), [spans, nodeType]);
  const filteredNodes = useMemo(() => TimelineService.getFilteredNodes(nodes, collapsedSpans), [collapsedSpans, nodes]);
  const [min, max] = useMemo(() => TimelineService.getMinMax(nodes), [nodes]);
  const getScale = useCallback(() => TimelineService.createScaleFunc({min, max}), [max, min]);

  const onSpanClick = useCallback(
    (spanId: string) => {
      TraceAnalyticsService.onTimelineSpanClick(spanId);
      onClick(spanId);
    },
    [onClick]
  );

  const onSpanNavigation = useCallback(
    (spanId: string) => {
      onNavigate(spanId);
      // TODO: Improve the method to search for the index
      const index = filteredNodes.findIndex(node => node.data.id === spanId);
      if (index !== -1) {
        listRef?.current?.scrollToItem(index, 'start');
      }
    },
    [filteredNodes, listRef, onNavigate]
  );

  const onSpanCollapse = useCallback((spanId: string) => {
    setCollapsedSpans(prevCollapsed => {
      if (prevCollapsed.includes(spanId)) {
        return without(prevCollapsed, spanId);
      }
      return [...prevCollapsed, spanId];
    });
  }, []);

  const value = useMemo<IContext>(
    () => ({
      collapsedSpans,
      getScale: getScale(),
      matchedSpans,
      nameColumnWidth,
      onNameColumnWidthChange: setNameColumnWidth,
      onSpanClick,
      onSpanCollapse,
      onSpanNavigation,
      selectedSpan,
      spans: filteredNodes,
      viewEnd: max,
      viewStart: min,
    }),
    [
      collapsedSpans,
      filteredNodes,
      getScale,
      matchedSpans,
      max,
      min,
      onSpanClick,
      onSpanCollapse,
      onSpanNavigation,
      selectedSpan,
      nameColumnWidth,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TimelineProvider;
