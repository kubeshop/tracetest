import {NodeTypesEnum} from 'constants/Visualization.constants';
import noop from 'lodash/noop';
import without from 'lodash/without';
import Span from 'models/Span.model';
import TimelineModel from 'models/Timeline.model';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {FixedSizeList as List} from 'react-window';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {selectSpan} from 'redux/slices/Trace.slice';
import TraceSelectors from 'selectors/Trace.selectors';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import TimelineService, {TScaleFunction} from 'services/Timeline.service';
import {TNode} from 'types/Timeline.types';

interface IContext {
  collapsedSpans: string[];
  getScale: TScaleFunction;
  matchedSpans: string[];
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
}

export const useTimeline = () => useContext(Context);

const TimelineProvider = ({children, listRef, nodeType, spans}: IProps) => {
  const dispatch = useAppDispatch();
  const [collapsedSpans, setCollapsedSpans] = useState<string[]>([]);
  const matchedSpans = useAppSelector(TraceSelectors.selectMatchedSpans);
  const selectedSpan = useAppSelector(TraceSelectors.selectSelectedSpan);

  const nodes = useMemo(() => TimelineModel(spans, nodeType), [spans, nodeType]);
  const filteredNodes = useMemo(() => TimelineService.getFilteredNodes(nodes, collapsedSpans), [collapsedSpans, nodes]);
  const [min, max] = useMemo(() => TimelineService.getMinMax(nodes), [nodes]);
  const getScale = useCallback(() => TimelineService.createScaleFunc({min, max}), [max, min]);

  const onSpanClick = useCallback(
    (spanId: string) => {
      TraceAnalyticsService.onTimelineSpanClick(spanId);
      dispatch(selectSpan({spanId}));
    },
    [dispatch]
  );

  const onSpanNavigation = useCallback(
    (spanId: string) => {
      dispatch(selectSpan({spanId}));
      // TODO: Improve the method to search for the index
      const index = filteredNodes.findIndex(node => node.data.id === spanId);
      if (index !== -1) {
        listRef?.current?.scrollToItem(index, 'start');
      }
    },
    [dispatch, filteredNodes, listRef]
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
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TimelineProvider;
