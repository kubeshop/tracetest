import {useCallback, useEffect} from 'react';
import {Node, NodeChange} from 'react-flow-renderer';

import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import SkeletonDiagram from 'components/SkeletonDiagram';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import DAG from 'components/Visualization/components/DAG';
import Timeline from 'components/Visualization/components/Timeline';
import {TestState} from 'constants/TestRun.constants';
import {useSpan} from 'providers/Span/Span.provider';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {initNodes, onNodesChange as onNodesChangeAction} from 'redux/slices/DAG.slice';
import DAGSelectors from 'selectors/DAG.selectors';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import {TSpan} from 'types/Span.types';
import {TTestRunState} from 'types/TestRun.types';

export interface IProps {
  runState: TTestRunState;
  spans: TSpan[];
  type: VisualizationType;
}

const Visualization = ({runState, spans, type}: IProps) => {
  const dispatch = useAppDispatch();
  const edges = useAppSelector(DAGSelectors.selectEdges);
  const nodes = useAppSelector(DAGSelectors.selectNodes);
  const {onSelectSpan, matchedSpans, onSetFocusedSpan, focusedSpan, selectedSpan} = useSpan();

  const {isOpen} = useTestSpecForm();

  useEffect(() => {
    dispatch(initNodes({spans}));
  }, [dispatch, spans]);

  useEffect(() => {
    if (selectedSpan) return;
    const firstSpan = spans.find(span => !span.parentId);
    onSelectSpan(firstSpan?.id ?? '');
  }, [onSelectSpan, selectedSpan, spans]);

  const onNodesChange = useCallback((changes: NodeChange[]) => dispatch(onNodesChangeAction({changes})), [dispatch]);

  const onNodeClick = useCallback(
    (event, {id}: Node) => {
      TraceDiagramAnalyticsService.onClickSpan(id);
      onSelectSpan(id);
    },
    [onSelectSpan]
  );

  const onNodeClickTimeline = useCallback(
    (spanId: string) => {
      TraceAnalyticsService.onTimelineSpanClick(spanId);
      onSelectSpan(spanId);
    },
    [onSelectSpan]
  );

  const onNavigateToSpan = useCallback(
    (spanId: string) => {
      onSelectSpan(spanId);
      onSetFocusedSpan(spanId);
    },
    [onSelectSpan, onSetFocusedSpan]
  );

  if (runState !== TestState.FINISHED) {
    return <SkeletonDiagram />;
  }

  return type === VisualizationType.Dag ? (
    <DAG
      edges={edges}
      isMatchedMode={matchedSpans.length > 0 || isOpen}
      matchedSpans={matchedSpans}
      nodes={nodes}
      onNavigateToSpan={onNavigateToSpan}
      onNodesChange={onNodesChange}
      onNodeClick={onNodeClick}
      selectedSpan={focusedSpan}
    />
  ) : (
    <Timeline
      isMatchedMode={matchedSpans.length > 0 || isOpen}
      matchedSpans={matchedSpans}
      onNavigateToSpan={onNavigateToSpan}
      onNodeClick={onNodeClickTimeline}
      selectedSpan={selectedSpan?.id ?? ''}
      spans={spans}
    />
  );
};

export default Visualization;
