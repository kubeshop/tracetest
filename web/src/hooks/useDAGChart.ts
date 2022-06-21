import {noop} from 'lodash';
import {useEffect, useState} from 'react';

import DAGService, {TNode} from 'services/DAG.service';
import {TSpan} from 'types/Span.types';

export const useDAGChart = <T>(items: TNode<T>[], selectedSpan?: TSpan, onSelectSpan = noop) => {
  const [nodes, setNodes] = useState([]);
  const [edges, setEdges] = useState([]);

  useEffect(() => {
    console.log('### useDAGChart: effect');
    if (!items.length) return;

    const {nodes: generatedNodes, edges: generatedEdges} = DAGService.getNodesAndEdges(items);
    setNodes(generatedNodes);
    setEdges(generatedEdges);

    if (!selectedSpan) {
      onSelectSpan(generatedNodes[0].id);
    }
  }, [items]);

  return {nodes, edges};
};
