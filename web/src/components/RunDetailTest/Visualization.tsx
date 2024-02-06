import {useCallback, useEffect} from 'react';
import {Node, NodeChange} from 'react-flow-renderer';

import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import RunEvents from 'components/RunEvents';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import DAG from 'components/Visualization/components/DAG';
import Timeline from 'components/Visualization/components/Timeline';
import {TestRunStage} from 'constants/TestRunEvents.constants';
import {NodeTypesEnum} from 'constants/Visualization.constants';
import Span from 'models/Span.model';
import TestRunEvent from 'models/TestRunEvent.model';
import {useSpan} from 'providers/Span/Span.provider';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {initNodes, onNodesChange as onNodesChangeAction} from 'redux/slices/DAG.slice';
import DAGSelectors from 'selectors/DAG.selectors';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import TestRunService from 'services/TestRun.service';
import {TTestRunState} from 'types/TestRun.types';

export interface IProps {
  isDAGDisabled: boolean;
  runEvents: TestRunEvent[];
  runState: TTestRunState;
  spans: Span[];
  type: VisualizationType;
}

const Visualization = ({isDAGDisabled, runEvents, runState, spans, type}: IProps) => {
  const dispatch = useAppDispatch();
  const edges = useAppSelector(DAGSelectors.selectEdges);
  const nodes = useAppSelector(DAGSelectors.selectNodes);
  const {onSelectSpan, matchedSpans, onSetFocusedSpan, focusedSpan, selectedSpan} = useSpan();

  const {isOpen} = useTestSpecForm();

  useEffect(() => {
    if (isDAGDisabled) return;
    dispatch(initNodes({spans}));
  }, [dispatch, isDAGDisabled, spans]);

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

  if (TestRunService.shouldDisplayTraceEvents(runState, spans.length)) {
    return <RunEvents events={runEvents} stage={TestRunStage.Trace} state={runState} />;
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
      nodeType={NodeTypesEnum.TestSpan}
      onNavigateToSpan={onNavigateToSpan}
      onNodeClick={onNodeClickTimeline}
      selectedSpan={selectedSpan?.id ?? ''}
      spans={spans}
    />
  );
};

export default Visualization;
