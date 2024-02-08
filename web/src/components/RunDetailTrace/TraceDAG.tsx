import {useAppDispatch, useAppSelector} from 'redux/hooks';
import TraceSelectors from 'selectors/Trace.selectors';
import {Node, NodeChange} from 'react-flow-renderer';
import {changeNodes, initNodes, selectSpan} from 'redux/slices/Trace.slice';
import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import {useCallback, useEffect} from 'react';
import Trace from 'models/Trace.model';
import DAG from '../Visualization/components/DAG';
import LoadingSpinner, {SpinnerContainer} from '../LoadingSpinner';

interface IProps {
  trace: Trace;
  onNavigateToSpan(spanId: string): void;
}

const TraceDAG = ({trace: {spans}, onNavigateToSpan}: IProps) => {
  const matchedSpans = useAppSelector(TraceSelectors.selectMatchedSpans);
  const selectedSpan = useAppSelector(TraceSelectors.selectSelectedSpan);
  const nodes = useAppSelector(TraceSelectors.selectNodes);
  const edges = useAppSelector(TraceSelectors.selectEdges);
  const isMatchedMode = Boolean(matchedSpans.length);
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(initNodes({spans}));
  }, [dispatch, spans]);

  const onNodesChange = useCallback((changes: NodeChange[]) => dispatch(changeNodes({changes})), [dispatch]);

  const onNodeClick = useCallback(
    (event: React.MouseEvent, {id}: Node) => {
      event.stopPropagation();
      TraceDiagramAnalyticsService.onClickSpan(id);
      dispatch(selectSpan({spanId: id}));
    },
    [dispatch]
  );

  if (spans.length && !nodes.length) {
    return (
      <SpinnerContainer>
        <LoadingSpinner />
      </SpinnerContainer>
    );
  }

  return (
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
  );
};

export default TraceDAG;
