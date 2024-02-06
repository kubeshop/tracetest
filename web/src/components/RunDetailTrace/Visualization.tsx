import RunEvents from 'components/RunEvents';
import {TestRunStage} from 'constants/TestRunEvents.constants';
import {NodeTypesEnum} from 'constants/Visualization.constants';
import TestRunEvent from 'models/TestRunEvent.model';
import {useCallback, useEffect} from 'react';
import {Node, NodeChange} from 'react-flow-renderer';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {changeNodes, initNodes, selectSpan} from 'redux/slices/Trace.slice';
import TraceSelectors from 'selectors/Trace.selectors';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import TestRunService from 'services/TestRun.service';
import {TTestRunState} from 'types/TestRun.types';
import Span from 'models/Span.model';
import DAG from '../Visualization/components/DAG';
import Timeline from '../Visualization/components/Timeline';
import {VisualizationType} from './RunDetailTrace';

interface IProps {
  isDAGDisabled: boolean;
  runEvents: TestRunEvent[];
  runState: TTestRunState;
  spans: Span[];
  type: VisualizationType;
}

const Visualization = ({isDAGDisabled, runEvents, runState, spans, type}: IProps) => {
  const dispatch = useAppDispatch();
  const edges = useAppSelector(TraceSelectors.selectEdges);
  const matchedSpans = useAppSelector(TraceSelectors.selectMatchedSpans);
  const nodes = useAppSelector(TraceSelectors.selectNodes);
  const selectedSpan = useAppSelector(TraceSelectors.selectSelectedSpan);
  const isMatchedMode = Boolean(matchedSpans.length);

  useEffect(() => {
    if (isDAGDisabled) return;
    dispatch(initNodes({spans}));
  }, [dispatch, isDAGDisabled, spans]);

  useEffect(() => {
    if (selectedSpan) return;
    const firstSpan = spans.find(span => !span.parentId);
    dispatch(selectSpan({spanId: firstSpan?.id ?? ''}));
  }, [dispatch, selectedSpan, spans]);

  const onNodesChange = useCallback((changes: NodeChange[]) => dispatch(changeNodes({changes})), [dispatch]);

  const onNodeClick = useCallback(
    (event: React.MouseEvent, {id}: Node) => {
      event.stopPropagation();
      TraceDiagramAnalyticsService.onClickSpan(id);
      dispatch(selectSpan({spanId: id}));
    },
    [dispatch]
  );

  const onNodeClickTimeline = useCallback(
    (spanId: string) => {
      TraceAnalyticsService.onTimelineSpanClick(spanId);
      dispatch(selectSpan({spanId}));
    },
    [dispatch]
  );

  const onNavigateToSpan = useCallback(
    (spanId: string) => {
      dispatch(selectSpan({spanId}));
    },
    [dispatch]
  );

  if (TestRunService.shouldDisplayTraceEvents(runState, spans.length)) {
    return <RunEvents events={runEvents} stage={TestRunStage.Trace} state={runState} />;
  }

  return type === VisualizationType.Dag ? (
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
  ) : (
    <Timeline
      isMatchedMode={isMatchedMode}
      matchedSpans={matchedSpans}
      nodeType={NodeTypesEnum.TraceSpan}
      onNavigateToSpan={onNavigateToSpan}
      onNodeClick={onNodeClickTimeline}
      selectedSpan={selectedSpan}
      spans={spans}
    />
  );
};

export default Visualization;
