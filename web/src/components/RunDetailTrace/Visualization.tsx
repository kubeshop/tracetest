import RunEvents from 'components/RunEvents';
import {TestRunStage} from 'constants/TestRunEvents.constants';
import {NodeTypesEnum} from 'constants/Visualization.constants';
import TestRunEvent from 'models/TestRunEvent.model';
import {useCallback, useEffect} from 'react';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {selectSpan} from 'redux/slices/Trace.slice';
import TraceSelectors from 'selectors/Trace.selectors';
import TestRunService from 'services/TestRun.service';
import Trace from 'models/Trace.model';
import {TTestRunState} from 'types/TestRun.types';
import TimelineV2 from 'components/Visualization/components/Timeline/TimelineV2';
import {VisualizationType} from './RunDetailTrace';
import TraceDAG from './TraceDAG';

interface IProps {
  isDAGDisabled: boolean;
  runEvents: TestRunEvent[];
  runState: TTestRunState;
  trace: Trace;
  type: VisualizationType;
}

const Visualization = ({isDAGDisabled, runEvents, runState, trace, trace: {spans, rootSpan}, type}: IProps) => {
  const dispatch = useAppDispatch();
  const selectedSpan = useAppSelector(TraceSelectors.selectSelectedSpan);

  useEffect(() => {
    if (selectedSpan) return;

    dispatch(selectSpan({spanId: rootSpan.id ?? ''}));
  }, [dispatch, rootSpan.id, selectedSpan, spans]);

  const onNavigateToSpan = useCallback(
    (spanId: string) => {
      dispatch(selectSpan({spanId}));
    },
    [dispatch]
  );

  if (TestRunService.shouldDisplayTraceEvents(runState, spans.length)) {
    return <RunEvents events={runEvents} stage={TestRunStage.Trace} state={runState} />;
  }

  return type === VisualizationType.Dag && !isDAGDisabled ? (
    <TraceDAG trace={trace} onNavigateToSpan={onNavigateToSpan} />
  ) : (
    <TimelineV2 nodeType={NodeTypesEnum.TraceSpan} spans={spans} />
  );
};

export default Visualization;
