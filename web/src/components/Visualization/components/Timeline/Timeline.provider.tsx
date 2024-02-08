import noop from 'lodash/noop';
import {createContext, useCallback, useContext, useMemo} from 'react';
import {FixedSizeList as List} from 'react-window';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {selectSpan} from 'redux/slices/Trace.slice';
import TraceSelectors from 'selectors/Trace.selectors';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import TimelineService, {TScaleFunction} from 'services/Timeline.service';
import {TNode} from 'types/Timeline.types';

interface IContext {
  getScale: TScaleFunction;
  matchedSpans: string[];
  onNavigateToSpan(spanId: string): void;
  onNodeClick(spanId: string): void;
  selectedSpan: string;
}

export const Context = createContext<IContext>({
  getScale: () => ({start: 0, end: 0}),
  matchedSpans: [],
  onNavigateToSpan: noop,
  onNodeClick: noop,
  selectedSpan: '',
});

interface IProps {
  children: React.ReactNode;
  listRef: React.RefObject<List>;
  nodes: TNode[];
}

export const useTimeline = () => useContext(Context);

const TimelineProvider = ({children, listRef, nodes}: IProps) => {
  const dispatch = useAppDispatch();
  const matchedSpans = useAppSelector(TraceSelectors.selectMatchedSpans);
  const selectedSpan = useAppSelector(TraceSelectors.selectSelectedSpan);

  const [min, max] = useMemo(() => TimelineService.getMinMax(nodes), [nodes]);
  const getScale = useCallback(() => TimelineService.createScaleFunc({min, max}), [max, min]);

  const onNodeClick = useCallback(
    (spanId: string) => {
      TraceAnalyticsService.onTimelineSpanClick(spanId);
      dispatch(selectSpan({spanId}));
    },
    [dispatch]
  );

  const onNavigateToSpan = useCallback(
    (spanId: string) => {
      dispatch(selectSpan({spanId}));
      // TODO: Improve the method to search for the index
      const index = nodes.findIndex(node => node.data.id === spanId);
      listRef?.current?.scrollToItem(index, 'start');
    },
    [dispatch, listRef, nodes]
  );

  const value = useMemo<IContext>(
    () => ({
      getScale: getScale(),
      matchedSpans,
      onNavigateToSpan,
      onNodeClick,
      selectedSpan,
    }),
    [getScale, matchedSpans, onNavigateToSpan, onNodeClick, selectedSpan]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TimelineProvider;
