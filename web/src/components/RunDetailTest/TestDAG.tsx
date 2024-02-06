import {useCallback, useEffect} from 'react';
import {Node, NodeChange} from 'react-flow-renderer';

import DAG from 'components/Visualization/components/DAG';
import {useSpan} from 'providers/Span/Span.provider';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {initNodes, onNodesChange as onNodesChangeAction} from 'redux/slices/DAG.slice';
import DAGSelectors from 'selectors/DAG.selectors';
import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import Trace from 'models/Trace.model';
import {useTestSpecForm} from '../TestSpecForm/TestSpecForm.provider';
import LoadingSpinner from '../LoadingSpinner';

export interface IProps {
  trace: Trace;
  onNavigateToSpan(spanId: string): void;
}

const TestDAG = ({trace: {spans}, onNavigateToSpan}: IProps) => {
  const dispatch = useAppDispatch();
  const edges = useAppSelector(DAGSelectors.selectEdges);
  const nodes = useAppSelector(DAGSelectors.selectNodes);
  const {onSelectSpan, matchedSpans, focusedSpan} = useSpan();
  const {isOpen} = useTestSpecForm();

  useEffect(() => {
    dispatch(initNodes({spans}));
  }, [dispatch, spans]);

  const onNodesChange = useCallback((changes: NodeChange[]) => dispatch(onNodesChangeAction({changes})), [dispatch]);

  const onNodeClick = useCallback(
    (event, {id}: Node) => {
      TraceDiagramAnalyticsService.onClickSpan(id);
      onSelectSpan(id);
    },
    [onSelectSpan]
  );

  if (spans.length && !nodes.length) {
    return <LoadingSpinner />;
  }

  return (
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
  );
};

export default TestDAG;
