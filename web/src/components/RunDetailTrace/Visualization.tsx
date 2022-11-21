import SkeletonDiagram from 'components/SkeletonDiagram';
import {TestState} from 'constants/TestRun.constants';
import {useCallback, useEffect} from 'react';
import {Node, NodeChange} from 'react-flow-renderer';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {changeNodes, initNodes, selectSpan} from 'redux/slices/Trace.slice';
import TraceSelectors from 'selectors/Trace.selectors';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import {TSpan} from 'types/Span.types';
import {TTestRunState} from 'types/TestRun.types';
import {useDrawer} from '../Drawer/Drawer';
import DAG from '../Visualization/components/DAG';
import {Flame} from '../Visualization/components/Flame/Flame';
import Timeline from '../Visualization/components/Timeline';
import {VisualizationType} from './RunDetailTrace';

interface IProps {
  runState: TTestRunState;
  spans: TSpan[];
  type: VisualizationType;
}

const Visualization = ({runState, spans, type}: IProps) => {
  const dispatch = useAppDispatch();
  const edges = useAppSelector(TraceSelectors.selectEdges);
  const matchedSpans = useAppSelector(TraceSelectors.selectMatchedSpans);
  const nodes = useAppSelector(TraceSelectors.selectNodes);
  const selectedSpan = useAppSelector(TraceSelectors.selectSelectedSpan);
  const isMatchedMode = Boolean(matchedSpans.length);
  const {openDrawer} = useDrawer();

  useEffect(() => {
    dispatch(initNodes({spans}));
  }, [dispatch, spans]);

  useEffect(() => {
    if (selectedSpan) return;
    const firstSpan = spans.find(span => !span.parentId);
    dispatch(selectSpan({spanId: firstSpan?.id ?? ''}));
  }, [dispatch, selectedSpan, spans]);

  const onNodesChange = useCallback((changes: NodeChange[]) => dispatch(changeNodes({changes})), [dispatch]);

  const onNodeClick = useCallback(
    (event, {id}: Node) => {
      TraceDiagramAnalyticsService.onClickSpan(id);
      dispatch(selectSpan({spanId: id}));
      openDrawer();
    },
    [dispatch, openDrawer]
  );

  const onNodeClickTimeline = useCallback(
    (spanId: string) => {
      TraceAnalyticsService.onTimelineSpanClick(spanId);
      dispatch(selectSpan({spanId}));
      openDrawer();
    },
    [dispatch, openDrawer]
  );

  const onNavigateToSpan = useCallback(
    (spanId: string) => {
      dispatch(selectSpan({spanId}));
    },
    [dispatch]
  );

  if (runState !== TestState.FINISHED) {
    return <SkeletonDiagram />;
  }

  return (
    <>
      {type === VisualizationType.Dag && (
        <DAG
          edges={edges}
          isMatchedMode={isMatchedMode}
          matchedSpans={matchedSpans}
          nodes={nodes}
          onNavigateToSpan={onNavigateToSpan}
          onNodesChange={onNodesChange}
          onNodeClick={onNodeClick}
          selectedSpan={selectedSpan}
        />
      )}
      {type === VisualizationType.Timeline && (
        <Timeline
          isMatchedMode={isMatchedMode}
          matchedSpans={matchedSpans}
          onNavigateToSpan={onNavigateToSpan}
          onNodeClick={onNodeClickTimeline}
          selectedSpan={selectedSpan}
          spans={spans}
        />
      )}
      {type === VisualizationType.Flame && (
        <Flame
          isTrace
          isMatchedMode={isMatchedMode}
          matchedSpans={matchedSpans}
          onNavigateToSpan={onNavigateToSpan}
          onNodeClick={onNodeClickTimeline}
          selectedSpan={selectedSpan}
          spans={spans}
        />
      )}
    </>
  );
};

export default Visualization;
