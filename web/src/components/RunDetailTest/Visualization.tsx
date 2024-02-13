import {useCallback, useEffect} from 'react';

import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import RunEvents from 'components/RunEvents';
import Timeline from 'components/Visualization/components/Timeline';
import {TestRunStage} from 'constants/TestRunEvents.constants';
import {NodeTypesEnum} from 'constants/Visualization.constants';
import TestRunEvent from 'models/TestRunEvent.model';
import {useSpan} from 'providers/Span/Span.provider';
import Trace from 'models/Trace.model';
import TestRunService from 'services/TestRun.service';
import {TTestRunState} from 'types/TestRun.types';
import TestDAG from './TestDAG';

export interface IProps {
  isDAGDisabled: boolean;
  runEvents: TestRunEvent[];
  runState: TTestRunState;
  type: VisualizationType;
  trace: Trace;
}

const Visualization = ({isDAGDisabled, runEvents, runState, trace, trace: {spans, rootSpan}, type}: IProps) => {
  const {onSelectSpan, matchedSpans, onSetFocusedSpan, selectedSpan} = useSpan();

  useEffect(() => {
    if (selectedSpan) return;
    onSelectSpan(rootSpan.id);
  }, [onSelectSpan, rootSpan, selectedSpan, spans]);

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

  return type === VisualizationType.Dag && !isDAGDisabled ? (
    <TestDAG trace={trace} onNavigateToSpan={onNavigateToSpan} />
  ) : (
    <Timeline
      matchedSpans={matchedSpans}
      nodeType={NodeTypesEnum.TestSpan}
      onNavigate={onNavigateToSpan}
      onClick={onSelectSpan}
      selectedSpan={selectedSpan?.id ?? ''}
      spans={spans}
    />
  );
};

export default Visualization;
