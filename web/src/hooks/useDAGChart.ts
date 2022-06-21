import {noop} from 'lodash';
import {useEffect, useState} from 'react';
import {Node, Edge} from 'react-flow-renderer';

import DAGService, {INodeItem} from 'services/DAG.service';
import {TSpan} from 'types/Span.types';

export const useDAGChart = <T>(
  nodeList: INodeItem<T>[],
  affectedList: string[],
  matchedList: string[],
  selectedSpan?: TSpan,
  onSelectSpan = noop
) => {
  const [nodes, setNodes] = useState<Node<T>[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  useEffect(() => {
    console.log('### useDAGChart: effect');
    if (!nodeList.length) return;

    const {edges: generatedEdges, nodes: generatedNodes} = DAGService.getNodesAndEdges(nodeList);
    setEdges(generatedEdges);
    setNodes(generatedNodes);

    if (!selectedSpan) {
      onSelectSpan(generatedNodes[0].id);
    }
  }, [nodeList]);

  useEffect(() => {
    console.log('### useDAGChart: effect for affectedList, matchedList');

    setNodes(nds =>
      nds.map(node => {
        const isAffected = affectedList.includes(node.id);
        const isMatched = matchedList.includes(node.id);
        return {...node, className: `${isAffected ? 'affected' : ''} ${isMatched ? 'matched' : ''}`};
      })
    );
  }, [affectedList, matchedList]);

  return {nodes, edges};
};
