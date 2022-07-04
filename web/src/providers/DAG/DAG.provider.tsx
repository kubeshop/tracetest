import {noop} from 'lodash';
import {createContext, MouseEvent, useCallback, useContext, useEffect, useMemo, useState} from 'react';
import {Edge, Node, NodeChange} from 'react-flow-renderer';

import {useSpan} from 'providers/Span/Span.provider';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {initNodes, onNodesChange as onNodesChangeAction} from 'redux/slices/DAG.slice';
import TraceDiagramAnalyticsService from 'services/Analytics/TraceDiagramAnalytics.service';
import DAGSelectors from 'selectors/DAG.selectors';
import {TSpan} from 'types/Span.types';

const {onClickSpan} = TraceDiagramAnalyticsService;

interface IContext {
  edges: Edge[];
  isMiniMapActive: boolean;
  nodes: Node[];
  onMiniMapToggle(): void;
  onNodesChange(changes: NodeChange[]): void;
  onNodeClick(event: MouseEvent, node: Node): void;
}

const DagContext = createContext<IContext>({
  edges: [],
  isMiniMapActive: false,
  nodes: [],
  onMiniMapToggle: noop,
  onNodesChange: noop,
  onNodeClick: noop,
});

const useDAG = () => useContext(DagContext);

interface IProps {
  children: React.ReactNode;
  spans: TSpan[];
}

const DAGProvider = ({children, spans}: IProps) => {
  const dispatch = useAppDispatch();
  const edges = useAppSelector(DAGSelectors.selectEdges);
  const nodes = useAppSelector(DAGSelectors.selectNodes);
  const [isMiniMapActive, setIsMiniMapActive] = useState(false);
  const {onSelectSpan} = useSpan();

  useEffect(() => {
    dispatch(initNodes({spans}));
    const firstSpan = spans.find(span => !span.parentId);
    onSelectSpan(firstSpan?.id ?? '');
  }, [dispatch, onSelectSpan, spans]);

  const onMiniMapToggle = useCallback(() => {
    setIsMiniMapActive(isActive => !isActive);
  }, []);

  const onNodesChange = useCallback((changes: NodeChange[]) => dispatch(onNodesChangeAction({changes})), [dispatch]);

  const onNodeClick = useCallback(
    (event, {id}: Node) => {
      onClickSpan(id);
      onSelectSpan(id);
    },
    [onSelectSpan]
  );

  const value = useMemo(
    () => ({
      edges,
      isMiniMapActive,
      nodes,
      onMiniMapToggle,
      onNodesChange,
      onNodeClick,
    }),
    [edges, isMiniMapActive, nodes, onMiniMapToggle, onNodesChange, onNodeClick]
  );

  return <DagContext.Provider value={value}>{children}</DagContext.Provider>;
};

export {DAGProvider, useDAG};
