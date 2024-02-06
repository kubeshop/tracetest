import {useCallback, useEffect} from 'react';

import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import RunEvents from 'components/RunEvents';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import Timeline from 'components/Visualization/components/Timeline';
import {TestRunStage} from 'constants/TestRunEvents.constants';
import {NodeTypesEnum} from 'constants/Visualization.constants';
import TestRunEvent from 'models/TestRunEvent.model';
import {useSpan} from 'providers/Span/Span.provider';
import {useAppDispatch} from 'redux/hooks';
import {initNodes} from 'redux/slices/DAG.slice';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import Trace from 'models/Trace.model';
import TestRunService from 'services/TestRun.service';
import {TTestRunState} from 'types/TestRun.types';
import TestDAG from './TestDAG';

export interface IProps {
  runEvents: TestRunEvent[];
  runState: TTestRunState;
  type: VisualizationType;
  trace: Trace;
}

const Visualization = ({runEvents, runState, trace, trace: {spans, traceId}, type}: IProps) => {
  const dispatch = useAppDispatch();
  const {onSelectSpan, matchedSpans, onSetFocusedSpan, selectedSpan} = useSpan();

  const {isOpen} = useTestSpecForm();

  useEffect(() => {
    dispatch(initNodes({spans}));
  }, [dispatch, spans, traceId]);

  useEffect(() => {
    if (selectedSpan) return;
    const firstSpan = spans.find(span => !span.parentId);
    onSelectSpan(firstSpan?.id ?? '');
  }, [onSelectSpan, selectedSpan, spans]);

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
    <TestDAG trace={trace} onNavigateToSpan={onNavigateToSpan} />
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
